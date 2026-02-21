export interface Confirmation {
    shown: boolean;
    persistent?: boolean;
    title: string;
    text: string;
    yesText?: string;
    onConfirm: () => void;
}

export const useConfirmationStore = defineStore('confirmations', () => {
    const confirmation = ref<Confirmation>({
        shown: false,
        title: '',
        text: '',
        persistent: false,
        onConfirm: () => {
        }
    });

    return {confirmation};
})
