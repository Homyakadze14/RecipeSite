import Skeleton from 'react-loading-skeleton';
import 'react-loading-skeleton/dist/skeleton.css';
import styles from './RecipeCardSkeletonLoader.module.scss';

export const RecipeCardSkeletonLoader = () => {
	return (
		<div className={styles.loader}>
			<Skeleton className={styles.image} height={250} width={250} />
			<div className={styles.infoContainer}>
				<Skeleton height={40} className={styles.title} />
				<div className={styles.tagsContainer}>
					<Skeleton height={18.5} width={100} className={styles.tag} />
					<Skeleton height={18.5} width={100} />
				</div>
				<Skeleton height={170} className={styles.description} />
			</div>
		</div>
	);
};
