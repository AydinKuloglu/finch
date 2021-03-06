package storage

import (
	"os"
	"testing"
	"time"

	"github.com/emreler/finch/config"
	"github.com/emreler/finch/models"
	"gopkg.in/mgo.v2/bson"
)

var s Storage
var userID string
var alertID string
var alert *models.Alert

func TestMain(m *testing.M) {
	config := config.NewConfig("../config.json")

	s = NewStorage(config.Mongo)
	os.Exit(m.Run())
}

func TestLogEvents(t *testing.T) {
	myAlertID := bson.NewObjectId()
	err := s.LogProcessAlert(&models.Alert{ID: myAlertID}, 200)

	if err != nil {
		t.Error(err)
	}

	err = s.LogProcessAlert(&models.Alert{ID: myAlertID}, 200)

	if err != nil {
		t.Error(err)
	}

	err = s.LogCreateAlert(&models.Alert{ID: bson.NewObjectId()})

	if err != nil {
		t.Error(err)
	}

	err = s.LogCreateUser(&models.User{ID: bson.NewObjectId()})

	if err != nil {
		t.Error(err)
	}

	res, err := s.GetAlertHistory(myAlertID.Hex(), 10)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}

func TestCreateUser(t *testing.T) {
	user := &models.User{Name: "foo", Email: "bar@usefinch.co"}

	err := s.CreateUser(user)

	if err != nil {
		t.Error(err)
		return
	}

	userID = user.ID.Hex()

	t.Logf("Created userID: %s", userID)
}

func TestCreateAlert(t *testing.T) {
	alert = models.NewAlert()

	alert.Name = "foo's alert"
	alert.User = bson.ObjectIdHex(userID)
	alert.AlertDate = time.Now().Add(10 * time.Second)
	alert.Data = "somedata"

	err := s.CreateAlert(alert)

	if err != nil {
		t.Error(err)
		return
	}

	alertID = alert.ID.Hex()

	t.Logf("Created alertID: %s", alertID)
}

func TestGetAlert(t *testing.T) {
	var err error

	alert, err = s.GetAlert(alertID)

	if err != nil {
		t.Error(err)
		return
	}

	if alert.Data != "somedata" {
		t.Errorf("Alert data is wrong")
		return
	}
}

func TestGetUserAlerts(t *testing.T) {
	alerts, err := s.GetUserAlerts(userID)

	if err != nil {
		t.Error(err)
		return
	}

	if len(alerts) == 1 && alerts[0].Data == "somedata" {
		return
	}

	t.Errorf("Invalid user alerts data")
}

func TestUpdateAlert(t *testing.T) {
	alert.Data = "updated"

	err := s.UpdateAlert(alert)

	if err != nil {
		t.Error(err)
	}

	testAlert, _ := s.GetAlert(alertID)

	if testAlert.Data != "updated" {
		t.Errorf("Updated alert has invalid data: %s", testAlert.Data)
	}
}
