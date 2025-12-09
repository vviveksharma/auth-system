import axios from "axios";

const API_URL = "http://localhost:8081";

export const GetProfile = async () => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/me`,
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

export const DeleteProfile = async () => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.delete(
      `${API_URL}/tenant`,
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

export const SetPassword = async () => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.delete(
      `${API_URL}/tenant`,
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