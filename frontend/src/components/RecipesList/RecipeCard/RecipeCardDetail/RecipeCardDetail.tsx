import { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { ChefIcon } from '../../../../assets/icons/ChefIcon';
import { LikeIcon } from '../../../../assets/icons/LikeIcon';
import { TimeIcon } from '../../../../assets/icons/TimeIcon';
import { useAuthStore } from '../../../../store/auth/useAuthStore';
import { useRecipesStore } from '../../../../store/recipes/useRecipesStore';
import { IRecipe } from '../../../../types/interfaces';
import { addNewLines } from '../../../../utils/utils';
import { Button } from '../../../Button/Button';
import { CommentsList } from '../../../CommentsList/CommentsList';
import { Complexity } from '../../../Complexity/Complexity';
import { TextField } from '../../../TextField/TextField';
import styles from './RecipeCardDetail.module.scss';

export const RecipeCardDetail = ({ recipe }: { recipe: IRecipe }) => {
	const isLoading = useRecipesStore(state => state.isLoading);
	const setIsLoading = useRecipesStore(state => state.setIsLoading);

	const isAuth = useAuthStore(state => state.isAuth);

	const isLiked = useRecipesStore(state => state.isLiked);
	const likes_count = useRecipesStore(state => state.likes_count);

	const likeRecipe = useRecipesStore(state => state.likeRecipe);
	const unlikeRecipe = useRecipesStore(state => state.unlikeRecipe);

	const getRecipeAndComments = useRecipesStore(
		state => state.getRecipeAndComments
	);

	const comments = useRecipesStore(state => state.comments);
	const createComment = useRecipesStore(state => state.createComment);

	const text = useRecipesStore(state => state.text);
	const setText = useRecipesStore(state => state.setText);

	const {
		id,
		title,
		about,
		photos_urls,
		author,
		complexity,
		need_time,
		ingridients,
		instructions,
	} = recipe;

	useEffect(() => {
		setIsLoading(true);
		id && getRecipeAndComments(id);
		setIsLoading(false);
	}, [id, isLiked, isLoading]);

	const handleLike = (id: number) => {
		isLiked ? unlikeRecipe(id) : likeRecipe(id);
	};

	return (
		<article className={styles.card}>
			<div>
				<span className={styles.title}>{title}</span>
				<div className={styles.imageContainer}>
					<img src={photos_urls} />
				</div>
				<div className={styles.bottomInfoContainer}>
					<div className={styles.authorContainer}>
						<span className={styles.author}>Автор:</span>
						<Link to={`/user/${author?.login}`}>{author?.login}</Link>
					</div>
					<div className={styles.complexityAndTime}>
						<div className={styles.complexityContainer}>
							<ChefIcon />
							<Complexity starsAmount={complexity} />
						</div>
						<div>
							<TimeIcon />
							<span className={styles.time}>{need_time}</span>
						</div>
					</div>
				</div>
			</div>
			<p className={styles.description}>{addNewLines(about)}</p>
			<ul className={styles.detailsList}>
				<li>
					<h3>Ингредиенты</h3>
					<p>{addNewLines(ingridients)}</p>
				</li>
				<li>
					<h3>Инструкция приготовления</h3>
					<p>{addNewLines(instructions)}</p>
				</li>
			</ul>
			{isAuth ? (
				<>
					<div className={styles.likeContainer}>
						<button
							className={styles.likeButton}
							onClick={() => handleLike(id)}
						>
							<LikeIcon fill={isLiked ? '#000' : 'none'} />
						</button>
						<span>{likes_count}</span>
					</div>
					<form>
						<TextField
							direction='column'
							label='Комментарии:'
							field='textarea'
							value={text}
							onChange={e => setText(e.target.value)}
							placeholder='Добавить комментарий'
						/>
						<Button
							disabled={!text}
							title={!text ? 'Текст комментария не может быть пустым' : ''}
							onClick={e => createComment(e, id, text)}
						>
							Отправить
						</Button>
					</form>
				</>
			) : (
				<p className={styles.needForAuthText}>
					Чтобы оставлять комментарии и ставить лайки, войдите в аккаунт
				</p>
			)}
			<CommentsList comments={comments} />
		</article>
	);
};
