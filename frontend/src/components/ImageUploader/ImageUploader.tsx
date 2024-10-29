import React, { useEffect, useState } from 'react';
import { IImageUploader } from '../../types/interfaces';
import styles from './ImageUploader.module.scss';

export const ImageUploader = ({
	onImageUpload,
	selectedImage,
	label,
}: IImageUploader) => {
	const [imagePreview, setImagePreview] = useState('');

	const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0] || null;
		console.log(file);

		if (file) {
			const reader = new FileReader();
			reader.onloadend = () => {
				setImagePreview(reader.result as string);
				onImageUpload(file, reader.result as string);
			};
			reader.readAsDataURL(file);
		} else {
			setImagePreview('');
			onImageUpload(null, null);
		}
	};

	useEffect(() => {
		if (selectedImage) {
			if (typeof selectedImage === 'string') {
				setImagePreview(selectedImage);
			} else {
				const reader = new FileReader();
				reader.onloadend = () => {
					setImagePreview(reader.result as string);
				};
				reader.readAsDataURL(selectedImage);
			}
		} else {
			setImagePreview('');
		}
	}, [selectedImage]);

	return (
		<div className={styles.container}>
			<span className={styles.label}>{label}</span>
			<input
				type='file'
				accept='image/*'
				onChange={handleImageChange}
				className={styles.input}
				id='image-upload'
			/>
			<label htmlFor='image-upload' className={styles.labelWrapper}>
				{imagePreview ? (
					<img
						src={imagePreview}
						alt='Preview'
						className={styles.imagePreview}
					/>
				) : (
					<div className={styles.uploadButton}>
						<span className={styles.uploadButtonText}>Загрузить</span>
					</div>
				)}
			</label>
		</div>
	);
};
