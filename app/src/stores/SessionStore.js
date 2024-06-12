import { EventEmitter } from 'events';

class SessionStore extends EventEmitter {
    constructor() {
        super();
        this.data = {}; // Você pode inicializar a Store com os dados iniciais, se necessário
    }

  
}

// Exemplo de uso da classe Store
const sessionStore = new SessionStore();
export default sessionStore;