package logging

import "testing"
import "errors"
import "os"
import "strconv"

var handler, err = GetBasicHandler("","")

func TestBasicHandler(t *testing.T) {
	handler, err := GetBasicHandler("","")
	if err != nil {
		t.Errorf("TestBasicHandler GetBasicHandler() returned %s", err)
	}
	err = handler.SetFormatString("%(dateTime) - [%(fileName) %(lineNo)] %(levelName) %(message)")
	if err != nil {
		t.Errorf("TestBasicHandler SetFormatString() returned %s", err)
	}
	err = handler.SetFilePath(".", "")
	if err != nil {
		t.Errorf("TestBasicHandler SetFilePath() returned %s", err)
	}
	handler.SetLogLevel(WARNING)
	log := GetLogger("TestBasicHandler")
	log.AddHandler(handler)
	// log.Debug("HELLO")
	// log.Warning("WARNING")
	// log.Error("ERROR")
}

func TestRotatingHandler(t *testing.T) {
	handler, err := GetRotatingHandler("","rotating.log")
	if err != nil {
		t.Errorf("TestRotatingHandler GetRotatingHandler() returned %s", err)
	}
	err = handler.SetFormatString("%(dateTime) - [%(fileName) %(lineNo)] %(levelName) %(message)")
	if err != nil {
		t.Errorf("TestRotatingHandler SetFormatString() returned %s", err)
	}
	err = handler.SetFilePath(".", "rotating.log")
	if err != nil {
		t.Errorf("TestRotatingHandler SetFilePath() returned %s", err)
	}
	err = handler.SetMaxFileSize(1000)
	if err != nil {
		t.Errorf("TestRotatingHandler SetMaxFileSize() returned %s", err)
	}
	err = handler.SetBackupCount(10)
	if err != nil {
		t.Errorf("TestRotatingHandler SetMaxFileSize() returned %s", err)
	}
	log := GetLogger("TestRotatingHandler")
	log.AddHandler(handler)
	for i := 1; i < 1000; i++ {
		log.Error("ERROR")
	}
	for i := 1; i <= 10; i++ {
		err = os.Remove("rotating.log." + strconv.Itoa(i))
		if err != nil {
			t.Errorf("TestRotatingHandler os.Remove() returned %s", err)
		}
	}
}

func TestSetFormatString(t *testing.T) {
	handler, err := GetBasicHandler("","")
	formatString := "%(dateTime"
	err = handler.SetFormatString(formatString)
	if err == nil {
		t.Errorf("TestSetFormatString SetFormatString() returned %s words, want %s", errors.New("error format \"%(dateTime\""), err)
	}
	formatString = "%(messages)"
	err = handler.SetFormatString(formatString)
	if err == nil {
		t.Errorf("TestSetFormatString SetFormatString() returned %s words, want %s", errors.New("error formatName %(messagea)"), err)
	}
	formatString = "%(name)-%(levelName)-%(pathName)-%(fileName)-%(funcName)-%(lineNo)-%(date)-%(unixTime)-%(dateTime)-%(weekday)-%(nanoSecond)-%(ascTime)-%(message)"
	err = handler.SetFormatString(formatString)
	if err != nil {
		t.Errorf("TestSetFormatString SetFormatString() returned %s words, want %s", err, nil)
	}
	log := GetLogger("TestSetFormatString")
	log.AddHandler(handler)
	log.Debug("HELLO")

}

func TestSetFilePath(t *testing.T) {
	handler, err := GetBasicHandler("","")
	err = handler.SetFilePath("aa", "TestSetFilePathlog")
	if err == nil {
		t.Errorf("TestSetFilePath SetFilePath() returned %s words, want %s", err, nil)
	}
	err = handler.SetFilePath(".", "TestSetFilePathlog.log")
	if err != nil {
		t.Errorf("TestSetFilePath SetFilePath() returned %s words, want %s", nil, "error")
	}
	os.Remove("TestSetFilePathlog.log")
}

func TestSetFormatter(t *testing.T) {
	handler, err := GetBasicHandler("","")
	if err != nil {
		t.Errorf("TestSetFormatter GetBasicHandler() returned %s", err)
	}
	handler.logConfig.formatString = "%(name)-%(levelName)-%(pathName)-%(fileName)-%(funcName)-%(lineNo)-%(date)-%(unixTime)-%(dateTime)-%(weekday)-%(nanoSecond)-%(ascTime)-%(message)"
	err = handler.setFormatter()
	if err != nil {
		t.Errorf("TestSetFormatter setFormatter() returned %s words, want %s", err, nil)
	}
	handler.logConfig.formatString = "%(dateTime"
	err = handler.setFormatter()
	if err == nil {
		t.Errorf("TestSetFormatter setFormatter() returned %s words, want %s", errors.New("error format \"%(dateTime\""), err)
	}
	handler.logConfig.formatString = "%(messages)"
	err = handler.setFormatter()
	if err == nil {
		t.Errorf("TestSetFormatter setFormatter() returned %s words, want %s", errors.New("error formatName %(messagea)"), err)
	}

}

func TestSetBackupCount(t *testing.T) {
	handler, err := GetRotatingHandler("","rotating.log")
	if err != nil {
		t.Errorf("TestSetBackupCount GetRotatingHandler() returned %s", err)
	}
	err = handler.SetBackupCount(-1)
	if err == nil {
		t.Errorf("TestSetBackupCount SetBackupCount() returned %s", err)
	}
	err = handler.SetBackupCount(0)
	if err != nil {
		t.Errorf("TestSetBackupCount SetBackupCount() returned %s", err)
	}
	err = handler.SetBackupCount(10)
	if err != nil {
		t.Errorf("TestSetBackupCount SetBackupCount() returned %s", err)
	}
}

func TestSetMaxFileSize(t *testing.T) {
	handler, err := GetRotatingHandler("","rotating.log")
	if err != nil {
		t.Errorf("TestSetMaxFileSize GetRotatingHandler() returned %s", err)
	}
	err = handler.SetMaxFileSize(-1)
	if err == nil {
		t.Errorf("TestSetMaxFileSize SetBackupCount() returned %s", err)
	}
	err = handler.SetMaxFileSize(0)
	if err == nil {
		t.Errorf("TestSetMaxFileSize SetBackupCount() returned %s", err)
	}
	err = handler.SetMaxFileSize(1 * 100 * 1024 * 1024)
	if err != nil {
		t.Errorf("TestSetMaxFileSize SetBackupCount() returned %s", err)
	}

}

func TestSetWhen(t *testing.T) {
	handler, err := GetTimeRotatingHandler("","timeRotating.log")
	if err != nil {
		t.Errorf("TestSetWhen GetRotatingHandler() returned %s", err)
	}
	err = handler.SetWhen("1999dddd")
	if err == nil {
		t.Errorf("TestSetWhen SetWhen() returned %s", err)
	}
	err = handler.SetWhen("d111d")
	if err == nil {
		t.Errorf("TestSetWhen SetWhen() returned %s", err)
	}
	err = handler.SetWhen("1111sd")
	if err == nil {
		t.Errorf("TestSetWhen SetWhen() returned %s", err)
	}
	err = handler.SetWhen("5s")
	if err != nil {
		t.Errorf("TestSetWhen SetWhen() returned %s", err)
	}
}
