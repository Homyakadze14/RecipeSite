export type Filter = {
	limit?: number;
	offset?: number;
	order_by?: -1 | 1;
	order_field?: string;
	query?: string;
};
