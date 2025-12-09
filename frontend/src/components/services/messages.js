import axios from "axios";

const API_URL = "http://localhost:8081";

export const ListMessages = async (page, page_size, status) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    // Build URL with optional status parameter
    let url = `${API_URL}/tenant/messages?page=${page}&page_size=${page_size}`;
    if (status) {
      url += `&status=${status}`;
    }

    const response = await axios.get(url, {
      headers: {
        Authorization: `Bearer ${authToken}`,
        "Content-Type": "application/json",
      },
    });

    return response.data;
  } catch (error) {
    console.error("Error while fetching messages:", error);
    throw error;
  }
};

export const ApproveMessage = async (id) => {
  try {
    const authToken = sessionStorage.getItem("authToken");
    console.log("the messageid : ", id);
    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.put(
      `${API_URL}/tenant/messages/approve?id=${id}`,
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
    console.error("Error while approving message:", error);
    throw error;
  }
};

export const RejectMessage = async (id) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.put(
      `${API_URL}/tenant/messages/reject?id=${id}`,
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
    console.error("Error while rejecting message:", error);
    throw error;
  }
};
