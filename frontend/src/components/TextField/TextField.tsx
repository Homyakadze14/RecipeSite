import { useState } from 'react';
import { HideIcon } from '../../assets/icons/HideIcon';
import { ShowIcon } from '../../assets/icons/ShowIcon';
import styles from './TextField.module.scss';
import { ITextField } from '../../types/interfaces';

export const TextField = ({
	direction,
	label,
	field,
	value,
	onChange,
	placeholder,
	type = 'text',
	className,
}: ITextField) => {
	const [isPasswordVisible, setIsPasswordVisible] = useState(false);

	return (
		<div className={[styles.textField, styles[direction], className].join(' ')}>
			<label className={styles.label}>{label}</label>
			{field === 'input' && (
				<input
					className={[
						styles.input,
						type === 'password' ? styles.passwordInput : '',
					].join(' ')}
					value={value}
					onChange={onChange}
					placeholder={placeholder}
					type={isPasswordVisible ? 'text' : type}
				/>
			)}
			{type === 'password' && field === 'input' && (
				<button
					className={styles.eyeButton}
					onClick={e => {
						e.preventDefault();
						setIsPasswordVisible(prev => !prev);
					}}
				>
					{isPasswordVisible ? <HideIcon /> : <ShowIcon />}
				</button>
			)}
			{field === 'textarea' && (
				<textarea
					className={styles.textarea}
					value={value}
					onChange={onChange}
					placeholder={placeholder}
				/>
			)}
		</div>
	);
};
