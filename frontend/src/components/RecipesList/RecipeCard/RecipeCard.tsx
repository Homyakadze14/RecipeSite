import { useNavigate } from 'react-router-dom';
import { ChefIcon } from '../../../assets/icons/ChefIcon';
import { DeleteIcon } from '../../../assets/icons/DeleteIcon';
import { EditIcon } from '../../../assets/icons/EditIcon';
import { TimeIcon } from '../../../assets/icons/TimeIcon';
import { useAuthStore } from '../../../store/auth/useAuthStore';
import { useRecipesStore } from '../../../store/recipes/useRecipesStore';
import { IRecipeCard } from '../../../types/interfaces';
import { addNewLines } from '../../../utils/utils';
import { Complexity } from '../../Complexity/Complexity';
import styles from './RecipeCard.module.scss';

export const RecipeCard = ({
	id,
	imageSrc,
	title,
	description,
	starsAmount,
	time,
	author,
	onClick,
	onDelete,
}: IRecipeCard) => {
	const login = useAuthStore(state => state.login);

	const setRecipeId = useRecipesStore(state => state.setRecipeId);

	const navigate = useNavigate();

	const handleEdit = (e: React.MouseEvent<HTMLButtonElement>) => {
		e.stopPropagation();
		localStorage.setItem('recipeId', String(id));
		setRecipeId(id);
		navigate(`/edit_recipe/${id}`);
	};

	const handleDelete = (e: React.MouseEvent<HTMLButtonElement>) => {
		e.stopPropagation();
		onDelete(login, id);
	};

	return (
		<li className={styles.card} onClick={onClick}>
			<img src={imageSrc} />
			<div className={styles.infoContainer}>
				<div className={styles.titleAndButtonsContainer}>
					<span className={styles.title}>{title}</span>
					{author.login === login && (
						<div className={styles.buttons}>
							<button onClick={handleEdit}>
								<EditIcon />
							</button>
							<button onClick={handleDelete}>
								<DeleteIcon />
							</button>
						</div>
					)}
				</div>
				<div>
					<div className={styles.timeContainer}>
						<TimeIcon />
						<span className={styles.time}>{time}</span>
					</div>
					<div className={styles.complexityContainer}>
						<ChefIcon />
						<Complexity starsAmount={starsAmount} />
					</div>
				</div>
				<p className={styles.description}>{addNewLines(description)}</p>
			</div>
		</li>
	);
};
