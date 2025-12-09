import axios from "axios";

const API_URL = "http://localhost:8081";

export const GetDashBoardDetails = async () => {
  try {
    const authToken = sessionStorage.getItem("authToken");
    const response = await axios.get(API_URL + "/tenant/dashboard", {
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
    });
    return response.data;
  } catch (error) {
    console.log("eror while creating the tenant for the user: ", error);
    throw error;
  }
};

export const GetActiveRoles = async (page, page_size) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/tokens/status?page=${page}&&page_size=${page_size}&&status=active`,
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