import axios from "axios";

const API_URL = "http://localhost:8081";

export const ListUsers = async (status, page, page_size) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/users?page=${page}&page_size=${page_size}&status=${status}`,
      {
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      }
    );

    return response.data;
  } catch (error) {
    console.error("Error while fetching tokens:", error);
    throw error;
  }
}