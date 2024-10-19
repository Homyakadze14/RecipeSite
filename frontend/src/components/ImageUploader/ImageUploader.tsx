import React, { useEffect, useState } from 'react';

interface IImageUploader {
	onImageUpload: (file: File | null, imageSrc: string | null) => void;
	selectedImage: File | null;

	label: string;
}

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
		<div
			style={{
				display: 'flex',
				flexDirection: 'row',
				alignItems: 'center',
				gap: 10,
				marginBottom: 20,
			}}
		>
			<span style={{ fontSize: 28 }}>{label}</span>
			<input
				type='file'
				accept='image/*'
				onChange={handleImageChange}
				style={{ display: 'none' }}
				id='image-upload'
			/>
			<label htmlFor='image-upload' style={{ cursor: 'pointer' }}>
				{imagePreview ? (
					<img
						src={imagePreview}
						alt='Preview'
						style={{
							width: '100px',
							height: '100px',
							objectFit: 'cover',
							borderRadius: 12,
						}}
					/>
				) : (
					<div
						style={{
							width: '100px',
							height: '100px',
							border: '1px solid #000',
							display: 'flex',
							alignItems: 'center',
							justifyContent: 'center',
							borderRadius: 12,
						}}
					>
						<span style={{ fontSize: 20, color: '#00000091' }}>Загрузить</span>
					</div>
				)}
			</label>
		</div>
	);
};
