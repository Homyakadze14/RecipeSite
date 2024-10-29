import { IModal } from '../../types/interfaces';
import styles from './Modal.module.scss';

export const Modal = ({
	isActive,
	setIsActive,
	children,
	...props
}: IModal) => {
	return (
		<div
			className={styles.modal + ' ' + (isActive ? styles.active : '')}
			onClick={() => setIsActive(prev => !prev)}
		>
			<div
				className={styles.modalContent}
				onClick={e => e.stopPropagation()}
				{...props}
			>
				<div className={styles.closeIcon} onClick={() => setIsActive(false)}>
					&#x2715;
				</div>
				{children}
			</div>
		</div>
	);
};
