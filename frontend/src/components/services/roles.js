import axios from "axios";

const API_URL = "http://localhost:8081";

export const ListRoles = async (status, page, page_size, roleType) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/roles?page=${page}&page_size=${page_size}&roletype=${roleType}&status=${status}`,
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

export const GetRolePermissions = async (roleId) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.get(
      `${API_URL}/tenant/roles/persmissions?roleId=${roleId}`,
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

export const AddRoles = async (name, displayName, description, permissions) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.post(
      `${API_URL}/tenant/roles/`,
      {
        name: name,
        display_name: displayName,
        description: description,
        permissions: permissions,
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

export const EnableRole = async(roleId) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.put(
      `${API_URL}/tenant/roles/enable?roleId=${roleId}`,
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
    console.error("Error while enabling the role:", error);
    throw error;
  }
}

export const DisableRole = async(roleId) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.put(
      `${API_URL}/tenant/roles/disable?roleId=${roleId}`,
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
    console.error("Error while disabling the role:", error);
    throw error;
  }
}

export const DeleteRole = async(roleId) => {
  try {
    const authToken = sessionStorage.getItem("authToken");

    if (!authToken) {
      throw new Error("No auth token found");
    }

    const response = await axios.delete(
      `${API_URL}/tenant/roles?roleId=${roleId}`,
      {
        headers: {
          Authorization: `Bearer ${authToken}`,
          "Content-Type": "application/json",
        },
      }
    );

    return response.data;
  } catch (error) {
    console.error("Error while deleting the role:", error);
    throw error;
  }
}