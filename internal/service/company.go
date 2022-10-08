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
	"entetry/gotest/internal/repository/postgre"
)

const (
	companyAlreadyHasALogoErr = "company already has a logo"
	fileSaveError             = "file save error"
	imageExt                  = ".jpeg"
	redisExp                  = 25 * time.Second
)

// Company service company struct
type Company struct {
	companyRepository postgre.CompanyRepository
	logoRepository    postgre.LogoRepository
	redis             *redis.Client
}

// NewCompany creates new Company service
func NewCompany(companyRepository postgre.CompanyRepository,
	logoRepository postgre.LogoRepository, redisClient *redis.Client) *Company {
	return &Company{
		companyRepository: companyRepository, logoRepository: logoRepository, redis: redisClient}
}

// GetAll return all companies
func (c *Company) GetAll(ctx context.Context) ([]*model.Company, error) {
	return c.companyRepository.GetAll(ctx)
}

// GetByID return company by its uuid
func (c *Company) GetByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
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

// Create create company
func (c *Company) Create(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	return c.companyRepository.Create(ctx, company)
}

// Update update company
func (c *Company) Update(ctx context.Context, company *model.Company) error {
	return c.companyRepository.Update(ctx, company)
}

// Delete delete company
func (c *Company) Delete(ctx context.Context, id uuid.UUID) error {
	return c.companyRepository.Delete(ctx, id)
}

// AddLogo add logo to a company( fails if company already has a logo)
func (c *Company) AddLogo(ctx context.Context, companyID string, file *multipart.FileHeader) error {
	id, err := uuid.Parse(companyID)
	if err != nil {
		return err
	}
	logo, err := c.logoRepository.GetByCompanyID(ctx, id)
	if err != nil {
		return err
	}
	if logo != nil {
		return fmt.Errorf(companyAlreadyHasALogoErr)
	}
	imageURI := c.buildFileURI(companyID)
	err = c.saveFile(imageURI, file)
	if err != nil {
		return fmt.Errorf(fileSaveError)
	}

	err = c.logoRepository.Create(ctx, id, imageURI)
	if err != nil {
		return err
	}
	return nil
}

// GetLogo Get company logo
func (c *Company) GetLogo(ctx context.Context, companyID uuid.UUID) (string, error) {
	logo, err := c.logoRepository.GetByCompanyID(ctx, companyID)
	if err != nil {
		return "", err
	}
	return logo.Image, nil
}

func (c *Company) buildFileURI(companyID string) string {
	wd, _ := os.Getwd()
	basepath := filepath.Join(wd, "data", "company")
	_ = os.MkdirAll(basepath, os.ModePerm)
	fileURI := filepath.Join(basepath, companyID)
	return fmt.Sprintf("%s%s", fileURI, imageExt)
}

func (c *Company) saveFile(fileName string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		if srcError := src.Close(); srcError != nil {
			log.Printf("Error closing file: %s\n", srcError)
		}
	}()

	dst, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		return err
	}
	defer func() {
		if dstError := dst.Close(); dstError != nil {
			log.Printf("Error closing file: %s\n", dstError)
		}
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	err = dst.Sync()
	if err != nil {
		return err
	}

	return nil
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
	err = c.redis.Set(ctx, company.ID.String(), b.Bytes(), redisExp).Err()
	if err != nil {
		log.Error(err)
	}
}
