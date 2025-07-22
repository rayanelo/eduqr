// Mock data for authentication
export const USERS = [
  {
    id: '8864c717-587d-472a-929a-8e5f298024da-0',
    displayName: 'Jaydon Frankie',
    email: 'demo@minimals.cc',
    password: 'demo1234',
    photoURL: '/assets/images/portraits/portrait_1.jpg',
    phoneNumber: '+40 777666555',
    country: 'Romania',
    address: '90210 Broadway Blvd, Nashville, TN 37011-5678, USA',
    state: 'Tennessee',
    city: 'Nashville',
    zipCode: '37011-5678',
    about: 'Praesent turpis. Phasellus viverra nulla ut metus varius laoreet. Phasellus tempus.',
    role: 'admin',
    isPublic: true,
  },
  {
    id: '8864c717-587d-472a-929a-8e5f298024da-1',
    displayName: 'John Doe',
    email: 'john@example.com',
    password: 'password123',
    photoURL: '/assets/images/portraits/portrait_2.jpg',
    phoneNumber: '+33 123456789',
    country: 'France',
    address: '123 Rue de la Paix, Paris, 75001, France',
    state: 'Île-de-France',
    city: 'Paris',
    zipCode: '75001',
    about: 'Développeur passionné par les nouvelles technologies.',
    role: 'user',
    isPublic: true,
  },
  {
    id: '8864c717-587d-472a-929a-8e5f298024da-2',
    displayName: 'Jane Smith',
    email: 'jane@example.com',
    password: 'password456',
    photoURL: '/assets/images/portraits/portrait_3.jpg',
    phoneNumber: '+1 5551234567',
    country: 'United States',
    address: '456 Main Street, New York, NY 10001, USA',
    state: 'New York',
    city: 'New York',
    zipCode: '10001',
    about: 'Designer créative avec une passion pour l\'expérience utilisateur.',
    role: 'user',
    isPublic: true,
  },
];

// Mock JWT tokens
const generateMockToken = (userId) => {
  const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const payload = btoa(JSON.stringify({
    sub: userId,
    iat: Math.floor(Date.now() / 1000),
    exp: Math.floor(Date.now() / 1000) + (24 * 60 * 60), // 24 hours
  }));
  const signature = btoa('mock-signature');
  return `${header}.${payload}.${signature}`;
};

// Mock authentication functions
export const mockAuth = {
  login: async (email, password) => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    const user = USERS.find(u => u.email === email && u.password === password);
    
    if (!user) {
      throw new Error('Invalid email or password');
    }
    
    const accessToken = generateMockToken(user.id);
    
    return {
      accessToken,
      user: {
        id: user.id,
        displayName: user.displayName,
        email: user.email,
        photoURL: user.photoURL,
        phoneNumber: user.phoneNumber,
        country: user.country,
        address: user.address,
        state: user.state,
        city: user.city,
        zipCode: user.zipCode,
        about: user.about,
        role: user.role,
        isPublic: user.isPublic,
      },
    };
  },
  
  register: async (email, password, firstName, lastName) => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // Check if user already exists
    const existingUser = USERS.find(u => u.email === email);
    if (existingUser) {
      throw new Error('User already exists');
    }
    
    const newUser = {
      id: `user-${Date.now()}`,
      displayName: `${firstName} ${lastName}`,
      email,
      password,
      photoURL: '/assets/images/portraits/portrait_default.jpg',
      phoneNumber: '',
      country: '',
      address: '',
      state: '',
      city: '',
      zipCode: '',
      about: '',
      role: 'user',
      isPublic: true,
    };
    
    // In a real app, you would save to database
    // USERS.push(newUser);
    
    const accessToken = generateMockToken(newUser.id);
    
    return {
      accessToken,
      user: {
        id: newUser.id,
        displayName: newUser.displayName,
        email: newUser.email,
        photoURL: newUser.photoURL,
        phoneNumber: newUser.phoneNumber,
        country: newUser.country,
        address: newUser.address,
        state: newUser.state,
        city: newUser.city,
        zipCode: newUser.zipCode,
        about: newUser.about,
        role: newUser.role,
        isPublic: newUser.isPublic,
      },
    };
  },
  
  getMyAccount: async (accessToken) => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // In a real app, you would decode the JWT token to get user ID
    // For mock purposes, we'll return the first user
    const user = USERS[0];
    
    return {
      user: {
        id: user.id,
        displayName: user.displayName,
        email: user.email,
        photoURL: user.photoURL,
        phoneNumber: user.phoneNumber,
        country: user.country,
        address: user.address,
        state: user.state,
        city: user.city,
        zipCode: user.zipCode,
        about: user.about,
        role: user.role,
        isPublic: user.isPublic,
      },
    };
  },
  
  isValidToken: (token) => {
    // Mock token validation
    return token && token.includes('.');
  },
}; 