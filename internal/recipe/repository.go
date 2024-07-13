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
		err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
			&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
			&recipe.PhotosUrls, &recipe.Created_at, &recipe.Updated_at)
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
	err := row.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
		&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
		&recipe.PhotosUrls, &recipe.Created_at, &recipe.Updated_at)

	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (repo *RecipeRepository) GetFullRecipe(ctx context.Context, id int) (*FullRecipe, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "SELECT * FROM recipes WHERE id=$1", id)
	recipe := &Recipe{}
	err = row.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
		&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
		&recipe.PhotosUrls, &recipe.Created_at, &recipe.Updated_at)
	if err != nil {
		return nil, err
	}

	fullRecipe := &FullRecipe{}
	row = tx.QueryRowContext(ctx, "SELECT login, icon_url FROM users WHERE id=$1", recipe.UserID)
	err = row.Scan(&fullRecipe.Author, &fullRecipe.AuthorIconUrl)
	if err != nil {
		return nil, err
	}
	fullRecipe.Recipe = recipe

	tx.Commit()

	return fullRecipe, nil
}

func (repo *RecipeRepository) Create(ctx context.Context, rp *Recipe) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO recipes(user_id,title,about,complexitiy,need_time,ingridients,photos_urls,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		rp.UserID, rp.Title, rp.About, rp.Complexitiy, rp.NeedTime, rp.Ingridients, rp.PhotosUrls, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (repo *RecipeRepository) Update(ctx context.Context, rp_id int, rp *Recipe) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE recipes SET title=$1,about=$2,complexitiy=$3,need_time=$4,ingridients=$5,photos_urls=$6,updated_at=$7 WHERE id=$8",
		rp.Title, rp.About, rp.Complexitiy, rp.NeedTime, rp.Ingridients, rp.PhotosUrls, time.Now(), rp_id)

	if err != nil {
		return err
	}

	return nil
}

func (repo *RecipeRepository) Delete(ctx context.Context, rp *Recipe) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM recipes WHERE id=$1", rp.ID)
	return err
}
