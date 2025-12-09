import axios from "axios";

const API_URL = "http://localhost:8081";

export const CreateTenant = async (email, password, name, organization) => {
  try {
    const response = await axios.post(API_URL + "/tenant", {
      name: name,
      email: email,
      password: password,
      campany: organization,
    });
    return response.data;
  } catch (error) {
    console.log("eror while creating the tenant for the user: ", error);
    throw error;
  }
};

export const LoginTenant = async (email, password) => {
  try {
    const response = await axios.post(API_URL + "/tenant/login", {
      email: email,
      password: password,
    });
    return response.data;
  } catch (error) {
    console.log("eror while logging in the tenant for the user: ", error);
    throw error;
  }
};
