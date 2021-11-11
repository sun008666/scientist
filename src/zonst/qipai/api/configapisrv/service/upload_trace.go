package service

import (
	"github.com/fwhezfwhez/errorx"
	"sync"
	"time"
	"zonst/qipai/api/configapisrv/dependency/db"
	"zonst/qipai/api/configapisrv/dependency/errs"
	"zonst/qipai/api/configapisrv/models"
)

type Trace struct {
	RequestId string
	m         *sync.Map
}

func NewTrace(requestId string) *Trace {
	var tmp = Trace{
		RequestId: requestId,
		m:         &sync.Map{},
	}
	return &tmp
}

func (t *Trace) StepStart(desc string) {
	var tmp = models.UploadTrace{
		RequestId:   t.RequestId,
		Description: desc,
	}
	if e := db.QipaiDB.Model(models.UploadTrace{}).Create(&tmp).Error; e != nil {
		errs.SaveError(errorx.Wrap(e))
		return
	}

	t.m.Store(tmp.Description, tmp)
}

func (t *Trace) StepOver(desc string) {
	tmpI, ok := t.m.Load(desc)
	if !ok {
		return
	}

	tmp := tmpI.(models.UploadTrace)

	if e := db.QipaiDB.Model(tmp).Where("id=?", tmp.Id).Updates(map[string]interface{}{
		"finish_at": time.Now(),
		"vstate":    2,
	}).Error; e != nil {
		errs.SaveError(errorx.Wrap(e))
		return
	}
}

func (t *Trace) StepFail(desc string, e error) {
	tmpI, ok := t.m.Load(desc)
	if !ok {
		return
	}

	tmp := tmpI.(models.UploadTrace)

	if e := db.QipaiDB.Model(tmp).Where("id=?", tmp.Id).Updates(map[string]interface{}{
		"finish_at":   time.Now(),
		"vstate":      3,
		"fail_reason": e.Error(),
	}).Error; e != nil {
		errs.SaveError(errorx.Wrap(e))
		return
	}
}
