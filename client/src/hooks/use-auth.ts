import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { User } from "@/types/api";
import api from "@/lib/axios-config";
import type { LoginFormData, SignupFormData } from "@/lib/validations";
import { type } from "os";
import { email } from "zod";

const getCookieValue = (name: string): string | null => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
  return null;
};

export const isAuthenticated = (): boolean => {
  return getCookieValue("logged_in") === "true";
};

// WARN: fix unmarshalling of type in user type in here
export const useAuthUser = () => {
  return useQuery<User, Error>({
    queryKey: ["auth", "user"],
    queryFn: async () => {
      const response = await api.get("/api/v1/users/me");
      return response.data.data.users;
    },
    enabled: isAuthenticated(),
    staleTime: 5 * 60 * 1000,
    retry: (failureCount, error) => {
      if (error instanceof Error && "status" in error) {
        const status = (error as Error & { status: number }).status;
        if (status === 401 || status === 403) {
          return false;
        }
      }
      return failureCount < 3;
    },
  });
};

export const useAuthLogin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["auth", "login"],
    mutationFn: async (data: LoginFormData) => {
      const response = await api.post("/api/v1/auth/login", data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["auth", "user"] });
    },
  });
};

export const useAuthLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationKey: ["auth", "logout"],
    mutationFn: async () => {
      await api.post("/api/v1/auth/logout");
    },
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: ["auth"] });
      queryClient.removeQueries({ queryKey: ["api-key"] });

      queryClient.clear();

      document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      document.cookie =
        "logged_in=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      document.cookie =
        "refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    },
    onError: () => {
      queryClient.removeQueries({ queryKey: ["auth"] });
      queryClient.removeQueries({ queryKey: ["api-key"] });
      queryClient.clear();

      document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      document.cookie =
        "logged_in=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      document.cookie =
        "refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    },
  });
};

export const useAuthSignup = () => {
  return useMutation({
    mutationKey: ["auth", "signup"],
    mutationFn: async (data: SignupFormData) => {
      const { confirmPassword, ...signupData } = data;
      const response = await api.post("/api/v1/auth/signup", signupData);
      return response.data;
    },
  });
};

export const useSendVerficationEmail = () => {
  return useMutation({
    mutationKey: ["email", "send", "verfication"],
    mutationFn: async (data: { email: string }) => {
      const response = await api.post("/api/v1/email-otp/send", {
        email: data.email,
        type: "email-verification",
      });
      return response.data;
    },
  });
};

export const useVerfiyEmailOtp = () => {
  return useMutation({
    mutationKey: ["email", "otp", "email-verification"],
    mutationFn: async (data: { email: string; otp: string }) => {
      const response = await api.post("/api/v1/email-otp/verify", {
        email: data.email,
        type: "email-verification",
        otp: data.otp,
      });
      return response.data;
    },
  });
};

export const useSendForgotPasswordEmail = () => {
  return useMutation({
    mutationKey: ["email", "send", "forgot"],
    mutationFn: async (email: string) => {
      const response = await api.post("/api/v1/email-otp/send", {
        email: email,
        type: "forget-password",
      });
      return response.data;
    },
  });
};

export const useVerifyForgotPasswordOtp = () => {
  return useMutation({
    mutationKey: ["email", "otp", "forgot-password"],
    mutationFn: async (data: { email: string; otp: string }) => {
      const response = await api.post("/api/v1/email-otp/send", {
        email: data.email,
        otp: data.otp,
        type: "forgot-password",
      });
      return response.data;
    },
  });
};
