import { useState } from 'react';
import { DeleteIcon } from '../../../assets/icons/DeleteIcon';
import { EditIcon } from '../../../assets/icons/EditIcon';
import { useAuthStore } from '../../../store/auth/useAuthStore';
import { useRecipesStore } from '../../../store/recipes/useRecipesStore';
import { IComment } from '../../../types/interfaces';
import { Button } from '../../Button/Button';
import { Modal } from '../../Modal/Modal';
import { TextField } from '../../TextField/TextField';
import styles from './Comment.module.scss';

export const Comment = ({ id, author, text }: IComment) => {
	const { icon_url } = author;

	const login = useAuthStore(state => state.login);

	const editComment = useRecipesStore(state => state.editComment);
	const deleteComment = useRecipesStore(state => state.deleteComment);

	const [isEditModalActive, setIsEditModalActive] = useState(false);
	const [editText, setEditText] = useState(text);

	const recipeId = Number(location.pathname.split('/').pop());

	return (
		<li className={styles.comment}>
			<div className={styles.info}>
				<div className={styles.user}>
					<img src={icon_url} alt='user avatar' />
					<span>{author.login}</span>
				</div>
				{login === author.login && (
					<div className={styles.buttons}>
						<button onClick={() => setIsEditModalActive(true)}>
							<EditIcon />
						</button>
						<button
							onClick={() => {
								deleteComment(id, recipeId);
								console.log('ids: ', id, recipeId);
							}}
						>
							<DeleteIcon />
						</button>
					</div>
				)}
			</div>
			<p className={styles.text}>{text}</p>

			<Modal isActive={isEditModalActive} setIsActive={setIsEditModalActive}>
				<div className={styles.modalContainer}>
					<form>
						<h2>Редактировать комментарий</h2>
						<TextField
							label='Текст:'
							direction='column'
							field='textarea'
							value={editText}
							onChange={e => setEditText(e.target.value)}
							placeholder='Текст комментария'
						/>
						<Button
							onClick={() => {
								editComment(id, editText);
								setIsEditModalActive(false);
							}}
						>
							Применить
						</Button>
					</form>
				</div>
			</Modal>
		</li>
	);
};
