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


type GetManyReferenceResponse<T> = {
    data: T[];
    total: number;
}

export const getManyReference = async <T>(resource: string, params?: Pagination): Promise<GetManyReferenceResponse<T>> => {
    let url = `${apiUrl}/${resource}`;
    if (params) {
        const { page, perPage } = params;

        const query = {
            _start: page * perPage,
            _end: (page + 1) * perPage,
            output_size_limit: 200,
        }

        url = `${apiUrl}/${resource}?${queryString.stringify(query)}`;
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