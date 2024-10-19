import { ChangeEvent, useState } from 'react';
import { HideIcon } from '../../assets/icons/HideIcon';
import { ShowIcon } from '../../assets/icons/ShowIcon';
import styles from './Input.module.scss';

export interface IInput {
	value: string;
	onChange: (event: ChangeEvent<HTMLInputElement>) => void;
	placeholder?: string;
	type?: string;
}

export const Input = ({
	value,
	onChange,
	placeholder,
	type = 'text',
}: IInput) => {
	const [isPasswordVisible, setIsPasswordVisible] = useState(false);

	return (
		<div className={styles.inputContainer}>
			<input
				value={value}
				onChange={onChange}
				placeholder={placeholder}
				className={styles.input}
				type={isPasswordVisible ? 'text' : type}
				style={{ paddingRight: type === 'password' ? 100 : 20 }}
			/>
			{type === 'password' && (
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
		</div>
	);
};
