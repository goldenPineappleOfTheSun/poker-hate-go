package players

type Personality struct {
	impulsiveness  float32
	bluffing       float32
	counting       float32
	intuition      float32
	aggressiveness float32
	/*normal_mood_chance float32
	  angry_mood_chance float32*/
}

func EmptyPersonality() Personality {
	return Personality{0, 0, 0, 0, 0}
}

func CreatePersonality(impulsiveness float32, bluffing float32, counting float32, intuition float32, aggressiveness float32) Personality {
	return Personality{impulsiveness, bluffing, counting, intuition, aggressiveness}
}

func (p *Personality) Multiply(other Personality) {
	p.impulsiveness = other.impulsiveness
	p.bluffing = other.bluffing
	p.counting = other.counting
	p.intuition = other.intuition
	p.aggressiveness = other.aggressiveness
}

/*
func (p *Personality) SetTransitionsChances(normal float32, angry float32) {
    p.normal_mood_chance = normal
    p.angry_mood_chance = angry
}

*/
