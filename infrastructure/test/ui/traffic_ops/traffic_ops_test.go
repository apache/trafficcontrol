package main

import (
	"testing"
	"time"

	"github.com/Comcast/traffic_control/infrastructure/test/environment"
	"github.com/tebeka/selenium"
)

func TestTrafficOpsLogin(t *testing.T) {
	env, err := environment.Get(environment.DefaultPath)
	if err != nil {
		t.Fatalf("Failed to get environment: %v\n", err)
	}

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		t.Fatalf("Error creating selenium Remote: %v\n", err)
	}
	defer wd.Quit()

	if err := wd.Get(env.TrafficOps.URI); err != nil {
		t.Fatalf("Error getting URI: %v\n", err)
	}

	elem, err := wd.FindElement(selenium.ByID, "u")
	if err != nil {
		t.Fatalf("Error finding element: %v\n", err)
	}
	elem.Clear()
	elem.SendKeys(env.TrafficOps.User)

	elem, err = wd.FindElement(selenium.ByID, "p")
	if err != nil {
		t.Fatalf("Error Finding element: %v\n", err)
	}

	elem.Clear()
	elem.SendKeys(env.TrafficOps.Password)

	btn, _ := wd.FindElement(selenium.ByID, "login_button")
	if err != nil {
		t.Fatalf("Error Finding element: %v\n", err)
	}
	btn.Click()

	loadingDiv, err := wd.FindElement(selenium.ByID, "utcclock")
	if err != nil {
		t.Fatalf("Error finding Element: %v\n", err)
	}
	for {
		if output, err := loadingDiv.Text(); err != nil {
			t.Fatalf("Error getting output: %v\n", err)
		} else if output != "loading..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	div, err := wd.FindElement(selenium.ByClassName, "dataTables_scroll")
	if err != nil {
		t.Fatalf("Error finding Element: %v\n", err)
	}

	//	txt, err := div.Text()
	_, err = div.Text()
	if err != nil {
		t.Fatalf("Error getting Text: %v\n", err)
	}
	//	fmt.Printf("Got: %s\n", txt)
}
