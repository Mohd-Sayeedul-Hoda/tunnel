import { useCallback, useState } from "react";
import { ZodError, ZodSchema } from "zod";

interface UseFormValidationOptions<T> {
  schema: ZodSchema<T>;
  validateOnBlur?: boolean;
  validateOnChange?: boolean;
}

interface ValidationState<T> {
  errors: Partial<Record<keyof T, string>>;
  touched: Partial<Record<keyof T, boolean>>;
  isValid: boolean;
}

export function useFormValidation<T extends Record<string, any>>({
  schema,
  validateOnBlur = true,
  validateOnChange = false,
}: UseFormValidationOptions<T>) {
  const [errors, setErrors] = useState<Partial<Record<keyof T, string>>>({});
  const [touched, setTouched] = useState<Partial<Record<keyof T, boolean>>>({});

  const validateField = useCallback(
    (field: keyof T, value: any, allData?: T) => {
      try {
        if (allData) {
          schema.parse(allData);
        } else {
          schema.pick({ [field]: true } as any).parse({ [field]: value });
        }
        return null;
      } catch (error) {
        if (error instanceof ZodError) {
          const fieldError = error.issues.find((err) => err.path[0] === field);
          return fieldError?.message || null;
        }
        return null;
      }
    },
    [schema],
  );

  const validateForm = useCallback(
    (data: T) => {
      try {
        schema.parse(data);
        return { isValid: true, errors: {} };
      } catch (error) {
        if (error instanceof ZodError) {
          const newErrors: Partial<Record<keyof T, string>> = {};
          error.issues.forEach((err) => {
            if (err.path[0]) {
              newErrors[err.path[0] as keyof T] = err.message;
            }
          });
          return { isValid: false, errors: newErrors };
        }
        return { isValid: false, errors: {} };
      }
    },
    [schema],
  );

  const handleFieldChange = useCallback(
    (field: keyof T, value: any, allData: T) => {
      if (!touched[field]) {
        setTouched((prev) => ({ ...prev, [field]: true }));
      }

      if (errors[field]) {
        setErrors((prev) => ({ ...prev, [field]: undefined }));
      }

      if (validateOnChange && touched[field]) {
        const fieldError = validateField(field, value, allData);
        if (fieldError) {
          setErrors((prev) => ({ ...prev, [field]: fieldError }));
        }
      }
    },
    [errors, touched, validateField, validateOnChange],
  );

  const handleFieldBlur = useCallback(
    (field: keyof T, value: any, allData: T) => {
      setTouched((prev) => ({ ...prev, [field]: true }));

      if (validateOnBlur) {
        const fieldError = validateField(field, value, allData);
        if (fieldError) {
          setErrors((prev) => ({ ...prev, [field]: fieldError }));
        }
      }
    },
    [validateField, validateOnBlur],
  );

  const setFieldError = useCallback((field: keyof T, error: string) => {
    setErrors((prev) => ({ ...prev, [field]: error }));
  }, []);

  const clearErrors = useCallback(() => {
    setErrors({});
  }, []);

  const reset = useCallback(() => {
    setErrors({});
    setTouched({});
  }, []);

  const isValid = Object.keys(errors).length === 0;

  return {
    errors,
    touched,
    isValid,
    validateField,
    validateForm,
    handleFieldChange,
    handleFieldBlur,
    setFieldError,
    clearErrors,
    reset,
  };
}
