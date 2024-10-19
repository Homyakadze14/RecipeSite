import styles from './Modal.module.scss';

export interface IModal extends React.HTMLAttributes<HTMLDivElement> {
	isActive: boolean;
	setIsActive: React.Dispatch<React.SetStateAction<boolean>>;
	children: React.ReactNode;
}

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
				{children}
			</div>
		</div>
	);
};
