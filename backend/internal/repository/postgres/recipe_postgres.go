package repo

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	"github.com/Homyakadze14/RecipeSite/internal/usecases"
	"github.com/Homyakadze14/RecipeSite/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

const (
	orderByAsc  = -1
	orderByDesc = 1
)

var constArraySize = 20

type RecipeRepo struct {
	*postgres.Postgres
}

func NewRecipeRepository(pg *postgres.Postgres) *RecipeRepo {
	return &RecipeRepo{pg}
}

func (r *RecipeRepo) GetAll(ctx context.Context) ([]entities.Recipe, error) {
	rows, err := r.Pool.Query(ctx, "SELECT * FROM recipes")

	if err != nil {
		return nil, fmt.Errorf("RecipeRepo - GetAll - r.Pool.Query: %w", err)
	}

	recipes := make([]entities.Recipe, 0, constArraySize)
	for rows.Next() {
		var recipe entities.Recipe
		err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
			&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
			&recipe.PhotosUrls, &recipe.CreatedAt, &recipe.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("RecipeRepo - GetAll - rows.Scan: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *RecipeRepo) GetFiltered(ctx context.Context, filter *entities.RecipeFilter) ([]entities.Recipe, error) {
	var request strings.Builder
	params := make([]interface{}, 0, 5)

	params = append(params, filter.Query)
	request.WriteString("SELECT * FROM recipes WHERE title LIKE '%'||$1||'%' OR about LIKE '%'||$1||'%' OR ingridients LIKE '%'||$1||'%'")

	allowOrderFields := []string{"", "title", "complexitiy", "updated_at"}
	if slices.Contains(allowOrderFields, filter.OrderField) {
		if filter.OrderField == "" {
			filter.OrderField = "title"
		}
		request.WriteString(fmt.Sprintf(" ORDER BY %s", filter.OrderField))
	} else {
		return nil, usecases.ErrBadOrderField
	}

	switch filter.OrderBy {
	case orderByAsc:
		request.WriteString(" ASC")
	case orderByDesc:
		request.WriteString(" DESC")
	}

	if filter.Limit != 0 {
		params = append(params, filter.Limit)
		request.WriteString(fmt.Sprintf(" LIMIT $%v", len(params)))
	}

	if filter.Offset != 0 {
		params = append(params, filter.Offset)
		request.WriteString(fmt.Sprintf(" OFFSET $%v", len(params)))
	}

	rows, err := r.Pool.Query(ctx, request.String(), params...)
	if err != nil {
		return nil, fmt.Errorf("RecipeRepo - GetFiltered - r.Pool.Query: %w", err)
	}

	recipes := make([]entities.Recipe, 0, constArraySize)
	for rows.Next() {
		var recipe entities.Recipe
		err := rows.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
			&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
			&recipe.PhotosUrls, &recipe.CreatedAt, &recipe.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("RecipeRepo - GetFiltered - rows.Scan: %w", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *RecipeRepo) Get(ctx context.Context, id int) (*entities.Recipe, error) {
	row := r.Pool.QueryRow(ctx, "SELECT * FROM recipes WHERE id=$1", id)

	recipe := &entities.Recipe{}
	err := row.Scan(&recipe.ID, &recipe.UserID, &recipe.Title, &recipe.About,
		&recipe.Complexitiy, &recipe.NeedTime, &recipe.Ingridients,
		&recipe.PhotosUrls, &recipe.CreatedAt, &recipe.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecases.ErrRecipeNotFound
		}
		return nil, fmt.Errorf("RecipeRepo - Get - row.Scan: %w", err)
	}

	return recipe, nil
}

func (r *RecipeRepo) Create(ctx context.Context, recipe *entities.Recipe) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO recipes(user_id,title,about,complexitiy,need_time,ingridients,photos_urls,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		recipe.UserID, recipe.Title, recipe.About, recipe.Complexitiy, recipe.NeedTime, recipe.Ingridients, recipe.PhotosUrls, time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("RecipeRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *RecipeRepo) Update(ctx context.Context, updatedRecipe *entities.Recipe) error {
	_, err := r.Pool.Exec(ctx, "UPDATE recipes SET title=$1,about=$2,complexitiy=$3,need_time=$4,ingridients=$5,photos_urls=$6,updated_at=$7 WHERE id=$8",
		updatedRecipe.Title, updatedRecipe.About, updatedRecipe.Complexitiy, updatedRecipe.NeedTime,
		updatedRecipe.Ingridients, updatedRecipe.PhotosUrls, time.Now(), updatedRecipe.ID)

	if err != nil {
		return fmt.Errorf("RecipeRepo - Update - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *RecipeRepo) Delete(ctx context.Context, recipe *entities.Recipe) error {
	_, err := r.Pool.Exec(ctx, "DELETE FROM recipes WHERE id=$1", recipe.ID)
	if err != nil {
		return fmt.Errorf("RecipeRepo - Delete - r.Pool.Exec: %w", err)
	}
	return nil
}
