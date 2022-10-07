package service

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository"
	"entetry/gotest/internal/repository/postgre"
)

const (
	CompanyAlreadyHasALogoErr = "company already has a logo"
	FileSaveError             = "file save error"
	ImageExt                  = ".jpeg"
)

type Company struct {
	companyRepository repository.CompanyRepository
	logoRepository    postgre.LogoRepository
	redis             *redis.Client
}

func NewCompany(companyRepository repository.CompanyRepository,
	logoRepository postgre.LogoRepository, redis *redis.Client) *Company {
	return &Company{
		companyRepository: companyRepository, logoRepository: logoRepository, redis: redis}
}

func (c *Company) GetAll(ctx context.Context) ([]*model.Company, error) {
	return c.companyRepository.GetAll(ctx)
}

func (c *Company) GetById(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	company, err := c.getCompanyRedis(ctx, id)
	if err != nil {
		log.Error(err)
	}
	if company != nil {
		return company, nil
	}
	company, err = c.companyRepository.GetOne(ctx, id)
	if company != nil {
		c.setCompanyRedis(ctx, company)
	}

	return company, err
}

func (c *Company) Create(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	return c.companyRepository.Create(ctx, company)
}

func (c *Company) Update(ctx context.Context, company *model.Company) error {
	return c.companyRepository.Update(ctx, company)
}

func (c *Company) Delete(ctx context.Context, id uuid.UUID) error {
	return c.companyRepository.Delete(ctx, id)
}

func (c *Company) AddLogo(ctx context.Context, companyId string, file *multipart.FileHeader) error {
	id, err := uuid.Parse(companyId)
	if err != nil {
		return err
	}
	logo, err := c.logoRepository.GetByCompanyID(ctx, id)
	if err != nil {
		return err
	}
	if logo != nil {
		return fmt.Errorf(CompanyAlreadyHasALogoErr)
	}
	imageUrl := c.buildFileUrl(companyId)
	err = c.saveFile(imageUrl, file)
	if err != nil {
		return fmt.Errorf(FileSaveError)
	}

	err = c.logoRepository.Create(ctx, id, imageUrl)
	if err != nil {
		return err
	}
	return nil
}

func (c *Company) buildFileUrl(companyId string) string {
	wd, _ := os.Getwd()
	basepath := filepath.Join(wd, "data", "company")
	os.MkdirAll(basepath, os.ModePerm)
	fileUri := filepath.Join(basepath, companyId)
	return fmt.Sprintf("%s%s", fileUri, ImageExt)
}

func (c *Company) saveFile(fileName string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	err = dst.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (c *Company) GetLogo(ctx context.Context, companyId uuid.UUID) (string, error) {
	logo, err := c.logoRepository.GetByCompanyID(ctx, companyId)
	if err != nil {
		return "", err
	}
	return logo.Image, nil
}

func (c *Company) getCompanyRedis(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	cmd := c.redis.Get(ctx, id.String())

	cmdb, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(cmdb)

	var res model.Company

	err = gob.NewDecoder(b).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Company) setCompanyRedis(ctx context.Context, company *model.Company) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(company)
	if err != nil {
		log.Error(err)
		return
	}
	err = c.redis.Set(ctx, company.ID.String(), b.Bytes(), 25*time.Second).Err()
	if err != nil {
		log.Error(err)
	}
}
