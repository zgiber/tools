package subscriptions

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

type Manager struct {
	sync.RWMutex
	topics *node
}

type Subscriber struct {
	w io.Writer
}

type node struct {
	next        map[string]*node       // next level
	empty       *node                  // the special case for empty topic segments. e.g.: "/" or "//" as topic
	mlwc        *node                  // multilevel wildcard
	slwc        *node                  // single level wildcard
	subscribers map[string]*subscriber // subscribers on this level (subtopic)
}

func newNode() *node {
	return &node{
		next:        map[string]*node{},
		subscribers: map[string]*subscriber{},
	}
}

func (sm *subscriptionManager) subscribe(s *subscriber, topic ...string) error {
	var err error
	for _, t := range topic {
		err = validateSubTopic(t)
		if err != nil {
			return err
		}

		// TODO: handle qos/will/retain
		if n := sm.topics.matchNode(splitTopic(t)); n != nil {
			n.subscribers[s.s.clientID] = s
		} else {
			n = sm.topics.createNode(splitTopic(t))
			n.subscribers[s.s.clientID] = s
		}
	}

	return nil
}

func (sm *subscriptionManager) unsubscribe(s *subscriber, topic ...string) error {
	var err error
	for _, t := range topic {
		err = validateSubTopic(t)
		if err != nil {
			return err
		}

		if n := sm.topics.matchNode(splitTopic(t)); n != nil {
			// TODO: qos/will/retain
			delete(n.subscribers, s.s.clientID)
		}
	}

	return nil
}

func (n *node) createNode(topic []string) *node {
	// direct match
	if len(topic) == 0 {
		return n
	}

	// next level topic
	switch topic[0] {
	case "":
		if n.empty == nil {
			n.empty = newNode()
			return n.empty
		}
		return n.empty.createNode(topic[1:])

	case "#":
		if n.mlwc == nil {
			n.mlwc = newNode()
			return n.mlwc
		}
		return n.mlwc.createNode(topic[1:])
	case "+":
		if n.slwc == nil {
			n.slwc = newNode()
			return n.slwc
		}
		return n.slwc.createNode(topic[1:])
	default:
		if next, ok := n.next[topic[0]]; ok {
			return next.createNode(topic[1:])
		} else {
			nn := newNode()
			n.next[topic[0]] = nn
			return nn.createNode(topic[1:])
		}
	}

}

func (n *node) matchNode(topic []string) *node {
	// direct match
	if len(topic) == 0 {
		return n
	}

	// next level topic
	switch topic[0] {
	case "":
		if n.empty == nil {
			//			n.empty = newNode()
			return n.empty
		}
		return n.empty.matchNode(topic[1:])

	case "#":
		if n.mlwc == nil {
			//			n.mlwc = newNode()
			return n.mlwc
		}
		return n.mlwc.matchNode(topic[1:])
	case "+":
		if n.slwc == nil {
			//			n.slwc = newNode()
			return n.slwc
		}
		return n.slwc.matchNode(topic[1:])
	default:
		if next, ok := n.next[topic[0]]; ok {
			return next.matchNode(topic[1:])
		}
		return nil
	}
}

// matchSubscribers takes a topic
func (n *node) matchSubscribers(topic []string, results *[]*subscriber) {

	// direct match
	if len(topic) == 0 {
		for _, s := range n.subscribers {
			*results = append(*results, s) // match direct subscribers
		}
		return
	}

	// next level topic
	switch topic[0] {
	case "":
		if n.empty != nil {
			n.empty.matchSubscribers(topic[1:], results)
		} else {
			n.empty = newNode()
		}

	case "#":

		for _, s := range n.subscribers {
			*results = append(*results, s)
		}

		if n.mlwc != nil {
			n.mlwc.matchSubscribers(topic[1:], results)
		} else {
			n.mlwc = newNode()
		}

		fallthrough

	case "+":

		for _, next := range n.next {
			next.matchSubscribers(topic[1:], results)
		}

		if n.slwc != nil {
			n.slwc.matchSubscribers(topic[1:], results)
		} else {
			n.slwc = newNode()
		}

		fallthrough

	default:
		if next, ok := n.next[topic[0]]; ok {
			next.matchSubscribers(topic[1:], results)
		}
	}
}

// "/hello" "hello" "/hello/" "hello/" are all valid, and represent diferent topics.
// " " is also a valid topic (since it is a unicode rune). A single "/" is also valid topic.
// "/" represents root
func splitTopic(topic string) []string {
	return strings.Split(topic, "/")
}

func validateSubTopic(topic string) error {
	levels := splitTopic(topic)

	for i := len(levels) - 1; i >= 0; i-- {
		switch levels[i] {
		case "#":
			if len(levels) != i+1 {
				return fmt.Errorf("Invalid topic. '#' must be last level.")
			}
			continue
		case "+":
			continue
		case "$":
			return fmt.Errorf("Invalid topic. '$' is not implemented yet.")

		}

		for j := len(levels[i]) - 1; j >= 0; j-- {
			switch levels[i][j] {
			case '#':
				return fmt.Errorf("Invalid topic. '#' must occupy entire topic level.")
			case '+':
				return fmt.Errorf("Invalid topic. '+' must occupy entire topic level.")
			case '$':
				return fmt.Errorf("Invalid topic. '$' is reserved for internal use.")
			case '\u0000':
				return fmt.Errorf("Invalid topic. Topic must not include null character.")
			}
		}
	}

	return nil
}

func validatePubTopic(topic string) error {
	if strings.ContainsAny(topic, "#+$\u0000") {
		return fmt.Errorf("Invalid topic.")
	}
	return nil
}
