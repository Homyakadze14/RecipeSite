package recipe

import (
	"context"
	"database/sql"
	"time"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{
		db: db,
	}
}

func (repo *RecipeRepository) GetAll(ctx context.Context) ([]Recipe, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM recipes")

	if err != nil {
		return nil, err
	}

	recipes := make([]Recipe, 0, 10)
	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.ID, &recipe.User_ID, &recipe.Title, &recipe.About,
			&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
			&recipe.Photos_URLS, &recipe.Created_at, &recipe.Updated_at)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (repo *RecipeRepository) Get(ctx context.Context, id int) (*Recipe, error) {
	row := repo.db.QueryRowContext(ctx, "SELECT * FROM recipes WHERE id=$1", id)

	recipe := &Recipe{}
	err := row.Scan(&recipe.ID, &recipe.User_ID, &recipe.Title, &recipe.About,
		&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
		&recipe.Photos_URLS, &recipe.Created_at, &recipe.Updated_at)

	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (repo *RecipeRepository) Create(ctx context.Context, rp *Recipe) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO recipes(user_id,title,about,complexitiy,need_time,ingridients,photos_urls,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		rp.User_ID, rp.Title, rp.About, rp.Complexitiy, rp.NeedTime, rp.Ingridients, rp.Photos_URLS, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (repo *RecipeRepository) Delete(ctx context.Context, rp *Recipe) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM recipes WHERE id=$1", rp.ID)
	return err
}
