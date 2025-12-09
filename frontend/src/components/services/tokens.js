import axios from "axios";

const API_URL = "http://localhost:8081";

export const CreateTokens = async (name, expiry) => {
  try {
    const authToken = sessionStorage.getItem("authToken");
    console.log(expiry);
    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.post(
      `${API_URL}/tenant/tokens`,
      {
        name: name,
        expiry_at: expiry,
      },
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
};

export const RevokeToken = async (token) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.put(
      `${API_URL}/tenant/tokens/${token}`,
      {},
      {
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      }
    );

    return response.data;
  } catch (error) {
    console.error("Error while revoking tokens:", error);
    throw error;
  }
};

export const ListTokensPaginted = async (page, page_size) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/tokens?page=${page}&&page_size=${page_size}`,
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
};

export const ListTokensWithStatus = async (status, page, page_size) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/tokens/status?page=${page}&&page_size=${page_size}&&status=${status}`,
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