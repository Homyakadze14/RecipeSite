import { IButton } from '../../types/interfaces';
import styles from './Button.module.scss';

export const Button = ({
	onClick,
	children,
	className,
	disabled = false,
	...props
}: IButton) => {
	return (
		<button
			className={styles.button + ' ' + className}
			style={{ cursor: disabled ? 'not-allowed' : 'pointer' }}
			disabled={disabled}
			onClick={onClick}
			{...props}
		>
			{children}
		</button>
	);
};
