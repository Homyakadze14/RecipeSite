import Skeleton from 'react-loading-skeleton';
import 'react-loading-skeleton/dist/skeleton.css';
import styles from './UserCardSkelteonLoader.module.scss';

export const UserCardSkeletonLoader = () => {
	return (
		<div className={styles.loader}>
			<div>
				<Skeleton className={styles.avatar} />
			</div>
			<div className={styles.infoContainer}>
				<div className={styles.nameAndButtonContainer}>
					<div>
						<Skeleton className={styles.name} />
					</div>
					<Skeleton className={styles.editButton} />
				</div>
				<Skeleton className={styles.description} />
			</div>
		</div>
	);
};
