package site

import (
	"context"
	"errors"

	"gorm.io/gorm"

	model "queueJob/pkg/db/table/job"
	"queueJob/pkg/tools/errs"
)

// 这种引入正常 是由于 go.mod 文件，module serverApi 确保模块声明与实际路径一致： "serverApi/pkg/tools/errs"

// CreatePunCompanyJob PunCompanyJob pun_company_job 企业招聘职位表
func (j *Job) CreatePunCompanyJob(ctx context.Context, model *model.PunCompanyJob) error {

	if j == nil || j.DB == nil {
		return errors.New("site.Site 或 DB 未初始化")
	}
	return j.DB.WithContext(ctx).Create(model).Error
	//return j.DB.WithContext(ctx).Create(model).Error
}

func (j *Job) DeletePunCompanyJob(ctx context.Context, id int) error {
	return j.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.PunCompanyJob{}).Error
}

func (j *Job) UpdatePunCompanyJob(ctx context.Context, id int, data map[string]interface{}) error {
	return j.DB.WithContext(ctx).Model(&model.PunCompanyJob{}).Where("id = ?", id).Updates(data).Error
}

func (j *Job) FindPunCompanyJobById(ctx context.Context, id int) (*model.PunCompanyJob, error) {
	var modelIn model.PunCompanyJob
	if err := j.DB.WithContext(ctx).First(&modelIn, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &modelIn, nil
}

func (j *Job) ListPunCompanyJob(ctx context.Context, offset, limit int) ([]*model.PunCompanyJob, error) {
	var models []*model.PunCompanyJob
	if err := j.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return models, nil
}

func (j *Job) CountPunCompanyJob(ctx context.Context) (int64, error) {
	var count int64
	if err := j.DB.WithContext(ctx).Model(&model.PunCompanyJob{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
