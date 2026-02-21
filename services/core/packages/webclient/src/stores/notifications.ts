export interface Notification {
    text: string;
    color?: string;
}

export const useNotificationStore = defineStore('notifications', () => {
    const queue = ref<Notification[]>([]);

    function add(message: Notification) {
        queue.value.push(message);
    }

    function info(message: string) {
        queue.value.push({text: message});
    }

    function error(message: string) {
        queue.value.push({text: message, color: 'error-container'});
    }

    return {queue, add, info, error}
})
