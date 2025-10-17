export interface CreateUserSkillInput {
  skill_id: number;
  user_id: number;
  proficiency_level?: 'beginner' | 'intermediate' | 'advanced' | 'expert';
}

export interface UpdateUserSkillInput {
  proficiency_level?: 'beginner' | 'intermediate' | 'advanced' | 'expert';
}

export interface UserSkillResponse {
  skill_id: number;
  user_id: number;
  proficiency_level: string;
  created_at: Date;
  skills?: {
    id: number;
    category: string;
    name: string;
    description?: string;
  };
  users?: {
    id: number;
    first_name: string;
    last_name?: string;
    login_email: string;
  };
}
