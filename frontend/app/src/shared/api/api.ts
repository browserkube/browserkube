import axios, { type AxiosRequestConfig, type AxiosError } from 'axios';
import { toast } from 'react-toastify';
import { BASE_URL } from '@shared/lib';

interface DefaultRequestProps {
  url: string;
  params?: AxiosRequestConfig['params'];
  allowErrorHandling?: boolean;
}

interface RequestPropsWithData extends DefaultRequestProps {
  data: AxiosRequestConfig['data'];
}

const getRequestUrl = (path: string) => {
  return `${BASE_URL}${path}`;
};

const getRequestConfig = (params?: AxiosRequestConfig['params']): AxiosRequestConfig => {
  const hasTimeoutInParams = !!params && !!params.timeout;
  return {
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: hasTimeoutInParams ? params.timeout : Number(process.env.REACT_APP_REQUEST_CONFIG_TIMEOUT),
    params,
  };
};

const getErrorMessage = (error: AxiosError, url: string) => {
  const errorRequest = url.replaceAll('/', '').toUpperCase();
  const errorCode = String(error.request.status) ?? '';
  const errorMessage = error.message ?? '';
  const errorCodeTemplate = errorCode ? `:${errorCode}` : '';
  const errorMessageTemplate = errorMessage ? ` : ${errorMessage}` : '';
  return `[${errorRequest}${errorCodeTemplate}]${errorMessageTemplate}`;
};

const errorHandler = (error: AxiosError, url: string) => {
  toast.error(getErrorMessage(error, url));
  if (process.env.NODE_ENV === 'development') {
    console.error(error);
  }
};

const get = async <Response>({ url, params, allowErrorHandling = true }: DefaultRequestProps): Promise<Response> => {
  const requestConfig = getRequestConfig(params);
  const response = await axios.get(getRequestUrl(url), requestConfig).catch(async (error) => {
    if (allowErrorHandling) {
      errorHandler(error, url);
    }
    return await Promise.reject(error);
  });
  return response.data;
};

const post = async <Response>({
  url,
  data,
  params,
  allowErrorHandling = true,
}: RequestPropsWithData): Promise<Response> => {
  const requestConfig = getRequestConfig(params);
  const response = await axios.post(getRequestUrl(url), data, requestConfig).catch(async (error) => {
    if (allowErrorHandling) {
      errorHandler(error, url);
    }
    return await Promise.reject(error);
  });
  return response.data;
};

const remove = async <Response = void>({
  url,
  params,
  allowErrorHandling = true,
}: DefaultRequestProps): Promise<Response> => {
  const requestConfig = getRequestConfig(params);
  const response = await axios.delete(getRequestUrl(url), requestConfig).catch(async (error) => {
    if (allowErrorHandling) {
      errorHandler(error, url);
    }
    return await Promise.reject(error);
  });
  return response.data;
};

const put = async <Response>({
  url,
  data,
  params,
  allowErrorHandling = true,
}: RequestPropsWithData): Promise<Response> => {
  const requestConfig = getRequestConfig(params);
  const response = await axios.put(getRequestUrl(url), data, requestConfig).catch(async (error) => {
    if (allowErrorHandling) {
      errorHandler(error, url);
    }
    return await Promise.reject(error);
  });
  return response.data;
};

export const api = { get, post, put, remove };
