-- Links each variant to its exact option values
CREATE TABLE variant_option_values (
    variant_id UUID NOT NULL REFERENCES variants(id) ON DELETE CASCADE,
    option_value_id UUID NOT NULL REFERENCES product_option_values(id) ON DELETE CASCADE,
    PRIMARY KEY (variant_id, option_value_id)
);

-- Indexes for fast joins
CREATE INDEX idx_variant_option_values_variant_id ON variant_option_values(variant_id);
CREATE INDEX idx_variant_option_values_option_value_id ON variant_option_values(option_value_id);
