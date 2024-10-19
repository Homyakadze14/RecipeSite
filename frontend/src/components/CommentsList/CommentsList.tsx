import { Comment, IComment } from './Comment/Comment';
import styles from './CommentsList.module.scss';

export interface ICommentsList {
	comments: IComment[];
}

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
