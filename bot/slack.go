package bots

import (
	"fmt"
	"log"

	client "github.com/nlopes/slack"
)

// slacktable struct
type slacktable struct {
	m     map[string]string
	color string
}

func (v *slacktable) Body() interface{} {
	params := client.PostMessageParameters{}
	attachment := client.Attachment{
		Fallback: fmt.Sprintf("%v", v.m),
		Fields:   []client.AttachmentField{},
	}
	for key, val := range v.m {
		attachment.Fields = append(attachment.Fields, client.AttachmentField{
			Title: key,
			Value: val,
			Short: true,
		})
	}
	attachment.Color = v.color
	params.Attachments = append(params.Attachments, attachment)
	return params
}

// SlackTable func
func SlackTable(color string, fs ...Metadata) Message {
	t := &slacktable{
		m:     make(map[string]string),
		color: color,
	}
	for _, md := range fs {
		t.m[md.Key()] = md.Val()
	}
	return t
}

// slack client
type slack struct {
	option    *Options
	client    *client.Client
	websocket *client.RTM
	receive   chan string
	queue     []Message
	connected bool
}

// Slack creates slack bot
func Slack(opts ...Option) (Bot, error) {
	v := new(slack)
	v.receive = make(chan string)
	v.option = newOptions(opts...)
	v.client = client.New(
		v.option.Token,
	)
	v.queue = make([]Message, 0)
	v.websocket = v.client.NewRTM()
	go v.websocket.ManageConnection()
	go v.listen()
	return v, nil
}

// listen to events
func (v *slack) listen() {
	for {
		select {
		case msg := <-v.websocket.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *client.HelloEvent:
			case *client.ConnectedEvent:
				v.connected = true
				if err := v.Emit(Text("Nap nap nap! I'm awake now ♥")); err != nil {
					log.Print(err)
				}
				for _, m := range v.queue {
					if err := v.Emit(m); err != nil {
						log.Print(err)
					}
				}
				v.queue = make([]Message, 0)
			case *client.DisconnectedEvent:
				v.connected = false
				if err := v.Emit(Text("Sorry gonna went for nap ♥")); err != nil {
					log.Print(err)
				}
			case *client.MessageEvent:
				v.receive <- ev.Text
			case *client.PresenceChangeEvent:
				log.Print("presence change:", ev)
			case *client.LatencyReport:
				log.Print("current latency:", ev.Value)
			case *client.RTMError:
				log.Print("rtm error", ev.Error())
			case *client.InvalidAuthEvent:
				log.Fatal("invalid credentials")
			}
		}
	}
}

func (v *slack) UserInfo() *client.User {
	user := v.websocket.GetInfo().User
	return v.websocket.GetInfo().GetUserByID(user.ID)
}

// Receive to handle receive message
func (v *slack) Receive() <-chan string {
	return v.receive
}

// Emit to send message
func (v *slack) Emit(m Message) error {
	if !v.connected {
		log.Print("slack connection is not established, message will be sent when connection is ready")
		v.queue = append(v.queue, m)
		return nil
	}
	bot := v.UserInfo()
	switch o := m.Body().(type) {
	case client.PostMessageParameters:
		o.Username = bot.Name
		o.IconURL = bot.Profile.ImageOriginal
		if _, _, err := v.client.PostMessage(v.option.Channel, o.Text, o); err != nil {
			return err
		}
	case string:
		if _, _, err := v.client.PostMessage(v.option.Channel, o, client.PostMessageParameters{
			Username: bot.Name,
			IconURL:  bot.Profile.ImageOriginal,
		}); err != nil {
			return err
		}
	default:
		log.Print("invalid message type")
	}
	return nil
}
