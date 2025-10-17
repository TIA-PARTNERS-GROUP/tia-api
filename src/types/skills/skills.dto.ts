export interface SkillResponse {
  id: number;
  category: string;
  name: string;
  description?: string;
  active: number;
  created_at: Date;
  user_skills_count?: number;
  project_skills_count?: number;
}

export interface SkillsFilter {
  category?: string;
  active?: boolean;
  search?: string;
}
