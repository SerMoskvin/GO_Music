// Базовые интерфейсы, расширяйте по мере необходимости
export interface User {
    user_id: number;
    login: string;
    role: string;
    surname: string;
    name: string;
    email: string;
}

export interface Section {
    name: string;
    url: string;
    can_read: boolean;
    can_write: boolean;
}

export interface RolePermissions {
    role: string;
    sections: Section[];
    own_records_only: boolean;
}

// Пример для сущности Employee
export interface Employee {
    employee_id: number;
    user_id?: number;
    surname: string;
    name: string;
    father_name?: string;
    birthday: string; // ISO string из JSON
    phone_number: string;
    job: string;
    work_experience: number;
}

// Обобщенный ответ API с пагинацией
export interface ApiResponse<T> {
    data: T[];
    total: number;
    page: number;
    page_size: number;
}