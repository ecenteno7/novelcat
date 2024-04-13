import axios from 'axios';

const BASE_URL = process.env.NODE_ENV === 'production' 
    ? 'https://www.dev.mybooks.tech' : 'http://localhost:8090';

async function seedDb() {
    try {
        const response = await axios.get(`${BASE_URL}/seedDb`);
        return response.data;
    } catch (error) {
        console.error('Error fetching books:', error);
        throw error;
    }
}

export { seedDb };