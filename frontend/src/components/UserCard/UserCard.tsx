import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { EditIcon } from '../../assets/icons/EditIcon';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useUsersStore } from '../../store/users/useUsersStore';
import { Button } from '../Button/Button';
import { ImageUploader } from '../ImageUploader/ImageUploader';
import { Modal } from '../Modal/Modal';
import { TextField } from '../TextField/TextField';
import styles from './UserCard.module.scss';
import { UserCardSkeletonLoader } from './UserCardSkeletonLoader/UserCardSkeletonLoader';

export interface IUserCard {
	avatarSrc: string;
	name: string;
	description: string;
}

export const UserCard = ({ avatarSrc, name, description }: IUserCard) => {
	const user = useUsersStore(state => state.user);
	const setUser = useUsersStore(state => state.setUser);

	const login = useAuthStore(state => state.login);

	const isAuth = useAuthStore(state => state.isAuth);

	const editUserForm = useUsersStore(state => state.editUserForm);
	const setEditUserForm = useUsersStore(state => state.setEditUserForm);

	const paramsLogin = useUsersStore(state => state.paramsLogin);

	const editUser = useUsersStore(state => state.editUser);

	const subscribe = useUsersStore(state => state.subscribe);
	const unsubscribe = useUsersStore(state => state.unsubscribe);

	const [isEditModalActive, setIsEditModalActive] = useState(false);

	const [isLoading, setIsLoading] = useState(true);

	const navigate = useNavigate();

	useEffect(() => console.log('params login: ', paramsLogin), [paramsLogin]);

	useEffect(() => {
		setTimeout(() => {
			setIsLoading(false);
		}, 200);
	}, []);

	useEffect(() => {
		setEditUserForm({
			login: user?.login,
			about: user?.about,
			icon: user?.icon_url,
		});
	}, [user?.login]);

	const handleSubscribe = () => {
		if (isAuth) {
			if (user.is_subscribed) {
				unsubscribe(paramsLogin);
				setUser({ ...user, is_subscribed: false });
			} else if (!user.is_subscribed) {
				subscribe(paramsLogin);
				setUser({ ...user, is_subscribed: true });
			}
		} else {
			alert('Ошибка! Чтобы подписаться, войдите в аккаунт.');
		}

		console.log('user is_subscribed: ', user.is_subscribed);
	};

	useEffect(() => console.log('editUserForm: ', editUserForm), [editUserForm]);

	const handleImageUpload = (image: File | null) => {
		setEditUserForm({
			...editUserForm,
			icon: image,
		});
	};

	return (
		<div className={styles.userCard}>
			{!isLoading && (
				<>
					<div>
						<img src={avatarSrc} />
					</div>
					<div className={styles.infoContainer}>
						<div className={styles.nameAndButtonContainer}>
							<div>
								<span>{name}</span>
								{name !== login && (
									<Button
										className={styles.subscribeButton}
										onClick={handleSubscribe}
									>
										{user.is_subscribed ? 'Отписаться' : 'Подписаться'}
									</Button>
								)}
							</div>
							{name === login && (
								<button
									className={styles.editButton}
									onClick={() => setIsEditModalActive(true)}
								>
									<EditIcon />
								</button>
							)}
						</div>
						<p>{description}</p>
					</div>
				</>
			)}
			{isLoading && <UserCardSkeletonLoader />}
			{isEditModalActive && (
				<Modal isActive={isEditModalActive} setIsActive={setIsEditModalActive}>
					<h3 style={{ textAlign: 'center', fontSize: '40px' }}>
						Изменить данные
					</h3>
					<form>
						<ImageUploader
							label='Аватар:'
							selectedImage={editUserForm.icon}
							onImageUpload={handleImageUpload}
						/>
						<TextField
							direction='row'
							label='Ваше имя:'
							field='input'
							value={editUserForm.login}
							placeholder='Введите ваше имя'
							onChange={e =>
								setEditUserForm({ ...editUserForm, login: e.target.value })
							}
						/>
						<TextField
							direction='column'
							label='Описание:'
							field='textarea'
							value={editUserForm.about}
							placeholder='Напишите что-нибудь о себе'
							onChange={e =>
								setEditUserForm({ ...editUserForm, about: e.target.value })
							}
						/>
						<div className={styles.changePasswordContainer}>
							<label className={styles.label} style={{ fontSize: '28px' }}>
								Пароль:
							</label>
							<Button onClick={() => navigate('/edit_password')}>
								Изменить пароль
							</Button>
						</div>

						<Button
							className={styles.saveButton}
							disabled={Object.entries(editUserForm).some(([key, value]) => {
								if (key !== 'about') !value;
							})}
							title={
								Object.entries(editUserForm).some(([key, value]) => {
									if (key !== 'about') !value;
								})
									? 'Имя и аватар должны быть заполнены'
									: ''
							}
							onClick={e => {
								setIsEditModalActive(false);
								editUser(
									e,
									login,
									editUserForm,
									navigate,
									setIsEditModalActive
								);
							}}
						>
							Применить
						</Button>
					</form>
				</Modal>
			)}
		</div>
	);
};