package database

import (
	"errors"
	"time"

	"gitlab-code-review-notifier/pkg/config"
)

var (
	ErrNotFound = errors.New("not found")
)

type ClientRepository struct {
	db *db
}

func NewClientRepository(db *db) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Get(id int) (*config.FiringConfig, error) {
	var clients []*config.FiringConfig
	err := r.db.Select(&clients, `select * from clients where id=$1`, id)
	if err != nil {
		return nil, err
	}

	if len(clients) == 0 {
		return nil, ErrNotFound
	}

	return clients[0], nil
}

func (r *ClientRepository) GetAll() ([]*config.FiringConfig, error) {
	clients := make([]*config.FiringConfig, 0)
	return clients, r.db.Select(&clients, `select * from clients`)
}

func (r *ClientRepository) Create(config *config.FiringConfig) error {
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	_, err := r.db.NamedExec(`insert into
			clients(
				group_id,
				gitlab_token,
				webhook_url,
				discussion_firing_timeout,
				merge_request_old_timeout,
				merge_request_old_mention,
				merge_request_review_timeout,
				merge_request_reviewers_count,
				merge_request_review_mention,
				created_at,
				updated_at
			)
			values (
				:group_id,
				:gitlab_token,
				:webhook_url,
				:discussion_firing_timeout,
				:merge_request_old_timeout,
				:merge_request_old_mention,
				:merge_request_review_timeout,
				:merge_request_reviewers_count,
				:merge_request_review_mention,
				:created_at,
				:updated_at
			)`,
		config)

	return err
}

func (r *ClientRepository) Update(config *config.FiringConfig) error {
	config.UpdatedAt = time.Now()

	_, err := r.db.NamedExec(`
			update clients set
				group_id=:group_id,
				gitlab_token=:gitlab_token,
				webhook_url=:webhook_url,
				discussion_firing_timeout=:discussion_firing_timeout,
				merge_request_old_timeout=:merge_request_old_timeout,
				merge_request_old_mention=:merge_request_old_mention,
				merge_request_review_timeout=:merge_request_review_timeout,
				merge_request_reviewers_count=:merge_request_reviewers_count,
				merge_request_review_mention=:merge_request_review_mention,
				updated_at=:updated_at
			where id=:id`,
		config)

	return err
}

func (r *ClientRepository) Delete(id int) error {
	_, err := r.db.Exec(`delete from clients where id=$1`, id)
	return err
}
