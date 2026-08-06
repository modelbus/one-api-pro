import { Message } from '@arco-design/web-vue'

export function toastSuccess(msg) { Message.success(msg) }
export function toastError(msg) { Message.error(msg) }
export function toastWarning(msg) { Message.warning(msg) }
export function toastInfo(msg) { Message.info(msg) }
