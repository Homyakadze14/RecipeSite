import { ICommentsList } from '../../types/interfaces';
import { Comment } from './Comment/Comment';
import styles from './CommentsList.module.scss';

export const CommentsList = ({ comments }: ICommentsList) => {
	return (
		<ul className={styles.list}>
			{comments.length ? (
				comments.map(comment => (
					<Comment
						key={comment.id}
						id={comment.id}
						author={comment.author}
						text={comment.text}
					/>
				))
			) : (
				<p className={styles.noDataText}>Ещё нет комментариев</p>
			)}
		</ul>
	);
};
