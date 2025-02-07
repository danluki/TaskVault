import axios from "axios";
import queryString from "query-string";

export const apiUrl = (window as any).TASKVAULT_API_URL || "http://localhost:8080/v1";

const httpClient = axios.create({
    baseURL: apiUrl,
    headers: {
        "Content-Type": "application/json",
    }
})

type Pagination = {
    page: number;
    perPage: number;
}

type Sort = {
    field: string;
    order: string;
}

type GetManyReferenceParams = {
    pagination: Pagination;
    sort: Sort;
    filter: Record<string, any>;
    target: string;
    id: string | number;
}

type GetManyReferenceResponse<T> = {
    data: T[];
    total: number;
}

export const getManyReference = async <T>(resource: string, params?: GetManyReferenceParams): Promise<GetManyReferenceResponse<T>> => {
    let url = `${apiUrl}/${resource}`;
    if (params) {
        const { page, perPage } = params.pagination;
        const { field, order }= params.sort;

        const query = {
            ...params.filter,
            [params.target]: params.id,
            _sort: field,
            _order: order,
            _start: (page - 1) * perPage,
            _end: page * perPage,
            output_size_limit: 200,
        }

        url = `${apiUrl}/${params.target}/${params.id}/${resource}?${queryString.stringify(query)}`;
    }

    try {
        const response = await httpClient.get<T[]>(url);
        const totalCount = response.headers["x-total-count"];
        if (!totalCount) {
            throw new Error(
                "The X-Total-Count header is missing in the HTTP response. Ensure it's included for proper pagination."
            );
        }

        return {
            data: response.data,
            total: parseInt(totalCount.split("/").pop() || "0", 10),
        };
    } catch (error) {
        console.error("Error fetching data:", error);
        throw error;
    }
}

export default httpClient;