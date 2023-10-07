package mail

import (
	"crypto/tls"
	"github.com/argcv/stork/config"
	"github.com/go-mail/mail"
	"io"
	"sync"
	"time"
)

// active fork
// https://github.com/go-mail/mail
type SMTPSession struct {
	Dialer *mail.Dialer
	Config *config.SMTPConfig

	isClosed   bool
	closeSch   *SMTPCloseSch
	sendCloser mail.SendCloser
}

func (s *SMTPSession) NewMessage() (*SMTPMessage) {
	m := &SMTPMessage{
		Dialer:  s.Dialer,
		Config:  s.Config,
		Message: mail.NewMessage(),

		locker: &sync.Mutex{},
		ss:     s,
	}
	m.DefaultFrom()
	return m
}

func (s *SMTPSession) Dial(f func(mail.SendCloser) error) error {
	s.closeSch.Activate()
	err := s.closeSch.SafeExec(func() error {
		if s.isClosed {
			//s.locker.Lock()
			//defer s.locker.Unlock()

			// dial new session
			sc, err := s.Dialer.Dial()

			if err != nil {
				return err
			}
			s.sendCloser = sc
			s.isClosed = false
		}
		return f(s.sendCloser)
	})
	return err
}

func (s *SMTPSession) close() error {
	if !s.isClosed {
		s.isClosed = true
		return s.sendCloser.Close()
	}
	return nil
}

func NewSMTPSession(c config.SMTPConfig) (s *SMTPSession, err error) {
	dialer := mail.NewDialer(c.Host, c.Port, c.User, c.Pass)
	if c.Option.InsecureSkipVerify {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	s = &SMTPSession{
		Dialer: dialer,
		Config: &c,

		isClosed: true,
	}
	s.closeSch = NewSMTPCloser(func() {
		s.close()
	})

	return
}

type SMTPMessage struct {
	Dialer  *mail.Dialer
	Message *mail.Message
	Config  *config.SMTPConfig

	locker *sync.Mutex
	ss     *SMTPSession
}

func (m *SMTPMessage) DefaultFrom() *SMTPMessage {
	if len(m.Config.Sender) > 0 {
		m.From(m.Config.User, m.Config.Sender)
	} else {
		m.From(m.Config.User)
	}
	return m
}

/** From Header can be used 1th.
 */
func (m *SMTPMessage) From(from string, name ...string) *SMTPMessage {
	if len(name) > 0 {
		m.SetHeader("From", m.FormatAddress(from, name[0]))
	} else {
		m.SetHeader("From", from)
	}
	return m
}

func (m *SMTPMessage) To(to string, name ...string) *SMTPMessage {
	if len(name) > 0 {
		m.AddToHeader("To", m.FormatAddress(to, name[0]))
	} else {
		m.AddToHeader("To", to)
	}
	return m
}

func (m *SMTPMessage) Cc(cc string, name ...string) *SMTPMessage {
	if len(name) > 0 {
		m.AddToHeader("Cc", m.FormatAddress(cc, name[0]))
	} else {
		m.AddToHeader("Cc", cc)
	}
	return m
}

func (m *SMTPMessage) Subject(subject string) *SMTPMessage {
	m.SetHeader("Subject", subject)
	return m
}

func (m *SMTPMessage) PlainBody(body string, alternative ...string) *SMTPMessage {
	m.SetBody("text/plain", body)
	if len(alternative) > 0 {
		m.AddAlternative("text/html", alternative[0])
	}
	return m
}

func (m *SMTPMessage) HtmlBody(body string) *SMTPMessage {
	m.SetBody("text/html", body)
	return m
}

func (m *SMTPMessage) Attach(path string, name ...string) *SMTPMessage {
	if len(name) > 0 {
		m.Message.Attach(path, mail.Rename(name[0]))
	} else {
		m.Message.Attach(path)
	}
	return m
}

func (m *SMTPMessage) AttachBytes(name string, data []byte) *SMTPMessage {
	m.Message.Attach(name, mail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(data)
		return err
	}))
	return m
}

func (m *SMTPMessage) WithDate() *SMTPMessage {
	m.SetHeader("X-Date", m.FormatDate(time.Now()))
	return m
}

func (m *SMTPMessage) SetHeader(field string, value ...string) *SMTPMessage {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.Message.SetHeader(field, value...)
	return m
}

func (m *SMTPMessage) AddToHeader(field string, value ...string) *SMTPMessage {
	m.locker.Lock()
	defer m.locker.Unlock()
	oh := m.Message.GetHeader(field)
	m.Message.SetHeaders(map[string][]string{
		field: append(oh, value...),
	})
	//locker.Message.SetHeader(field, append(locker.Message.GetHeader(field), value...)...)
	return m
}

func (m *SMTPMessage) AddAlternative(contentType, body string, settings ...mail.PartSetting) *SMTPMessage {
	m.Message.AddAlternative(contentType, body, settings...)
	return m
}

func (m *SMTPMessage) FormatDate(date time.Time) string {
	return m.Message.FormatDate(date)
}

func (m *SMTPMessage) FormatAddress(address, name string) string {
	return m.Message.FormatAddress(address, name)
}

//func (locker *SMTPMessage) SetHeaders(h map[string][]string) *SMTPMessage {
//	locker.Message.SetHeaders(h)
//	return locker
//}

func (m *SMTPMessage) SetBody(contentType, body string, settings ...mail.PartSetting) *SMTPMessage {
	m.Message.SetBody(contentType, body, settings...)
	return m
}

func (m *SMTPMessage) Perform() error {
	return m.ss.Dial(func(s mail.SendCloser) (err error) {
		err = mail.Send(s, m.Message)
		if err != nil {
			m.ss.close()
		}
		return
	})
	//return locker.Dialer.DialAndSend(locker.Message)
}

func (m *SMTPMessage) Send(to, subject, body string) error {
	m.To(to)
	m.Subject(subject)
	m.HtmlBody(body)
	return m.Perform()
}
