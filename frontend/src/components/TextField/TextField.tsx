import { useState } from 'react';
import { HideIcon } from '../../assets/icons/HideIcon';
import { ShowIcon } from '../../assets/icons/ShowIcon';
import styles from './TextField.module.scss';

export interface ITextField {
	direction: 'row' | 'column';
	label: string;
	field: 'input' | 'textarea';
	value: string;
	onChange: (
		event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => void;
	placeholder?: string;
	type?: string;
}

export const TextField = ({
	direction,
	label,
	field,
	value,
	onChange,
	placeholder,
	type = 'text',
}: ITextField) => {
	const [isPasswordVisible, setIsPasswordVisible] = useState(false);

	return (
		<div className={[styles.textField, styles[direction]].join(' ')}>
			<label className={styles.label}>{label}</label>
			{field === 'input' && (
				<input
					className={styles.input}
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
