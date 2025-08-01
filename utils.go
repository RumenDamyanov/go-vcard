package vcard

import (
	"fmt"
	"strings"
)

// escapeValue escapes special characters in vCard property values
func escapeValue(value string) string {
	// Replace special characters according to vCard specification
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, ",", "\\,")
	value = strings.ReplaceAll(value, ";", "\\;")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	return value
}

// unescapeValue unescapes special characters in vCard property values
func unescapeValue(value string) string {
	value = strings.ReplaceAll(value, "\\\\", "\\")
	value = strings.ReplaceAll(value, "\\,", ",")
	value = strings.ReplaceAll(value, "\\;", ";")
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\r", "\r")
	value = strings.ReplaceAll(value, "\\t", "\t")
	return value
}

// foldLine folds long lines according to vCard specification (75 characters)
func foldLine(line string) string {
	if len(line) <= 75 {
		return line
	}

	var result strings.Builder
	for i, r := range line {
		if i > 0 && i%75 == 0 {
			result.WriteString("\r\n ")
		}
		result.WriteRune(r)
	}

	return result.String()
}

// formatTypeParameter formats type parameters for vCard properties
func formatTypeParameter(types ...string) string {
	if len(types) == 0 {
		return ""
	}

	var validTypes []string
	for _, t := range types {
		if t != "" {
			validTypes = append(validTypes, t)
		}
	}

	if len(validTypes) == 0 {
		return ""
	}

	return ";TYPE=" + strings.Join(validTypes, ",")
}

// writeNameProperties writes name-related properties to the builder
func (v *VCard) writeNameProperties(builder *strings.Builder) error {
	// Write structured name (N property) - required
	builder.WriteString(fmt.Sprintf("N:%s\n", v.name.StructuredName()))

	// Write formatted name (FN property) - required
	formattedName := v.name.FormattedName()
	if formattedName == "" {
		// If no formatted name, use "Last, First" or just "First" or "Last"
		if v.name.Last != "" && v.name.First != "" {
			formattedName = v.name.Last + ", " + v.name.First
		} else if v.name.First != "" {
			formattedName = v.name.First
		} else if v.name.Last != "" {
			formattedName = v.name.Last
		}
	}

	if formattedName != "" {
		builder.WriteString(fmt.Sprintf("FN:%s\n", escapeValue(formattedName)))
	}

	return nil
}

// writeEmailProperties writes email properties to the builder
func (v *VCard) writeEmailProperties(builder *strings.Builder) {
	for _, email := range v.emails {
		var typeParam string
		if email.Type != "" {
			typeParam = formatTypeParameter(string(email.Type))
		} else {
			typeParam = formatTypeParameter("INTERNET")
		}

		if email.Preferred {
			typeParam += ";PREF=1"
		}

		line := fmt.Sprintf("EMAIL%s:%s", typeParam, escapeValue(email.Address))
		builder.WriteString(foldLine(line) + "\n")
	}
}

// writePhoneProperties writes phone properties to the builder
func (v *VCard) writePhoneProperties(builder *strings.Builder) {
	for _, phone := range v.phones {
		var typeParam string
		if phone.Type != "" {
			typeParam = formatTypeParameter(string(phone.Type))
		} else {
			typeParam = formatTypeParameter("VOICE")
		}

		if phone.Preferred {
			typeParam += ";PREF=1"
		}

		line := fmt.Sprintf("TEL%s:%s", typeParam, escapeValue(phone.Number))
		builder.WriteString(foldLine(line) + "\n")
	}
}

// writeAddressProperties writes address properties to the builder
func (v *VCard) writeAddressProperties(builder *strings.Builder) {
	for _, addr := range v.addresses {
		var typeParam string
		if addr.Type != "" {
			typeParam = formatTypeParameter(string(addr.Type))
		}

		if addr.Preferred {
			typeParam += ";PREF=1"
		}

		line := fmt.Sprintf("ADR%s:%s", typeParam, addr.StructuredAddress())
		builder.WriteString(foldLine(line) + "\n")

		// Also write formatted address label if we have address data
		if addr.Street != "" || addr.City != "" || addr.State != "" || addr.PostalCode != "" || addr.Country != "" {
			labelLine := fmt.Sprintf("LABEL%s:%s", typeParam, escapeValue(addr.FormattedAddress()))
			builder.WriteString(foldLine(labelLine) + "\n")
		}
	}
}

// writeOrganizationProperties writes organization properties to the builder
func (v *VCard) writeOrganizationProperties(builder *strings.Builder) {
	if v.organization.Name != "" {
		var orgParts []string
		orgParts = append(orgParts, escapeValue(v.organization.Name))
		if v.organization.Department != "" {
			orgParts = append(orgParts, escapeValue(v.organization.Department))
		}

		line := fmt.Sprintf("ORG:%s", strings.Join(orgParts, ";"))
		builder.WriteString(foldLine(line) + "\n")
	}

	if v.organization.Title != "" {
		line := fmt.Sprintf("TITLE:%s", escapeValue(v.organization.Title))
		builder.WriteString(foldLine(line) + "\n")
	}

	if v.organization.Role != "" {
		line := fmt.Sprintf("ROLE:%s", escapeValue(v.organization.Role))
		builder.WriteString(foldLine(line) + "\n")
	}
}

// writeURLProperties writes URL properties to the builder
func (v *VCard) writeURLProperties(builder *strings.Builder) {
	for _, url := range v.urls {
		var typeParam string
		if url.Type != "" {
			typeParam = formatTypeParameter(string(url.Type))
		}

		if url.Preferred {
			typeParam += ";PREF=1"
		}

		line := fmt.Sprintf("URL%s:%s", typeParam, escapeValue(url.Address))
		builder.WriteString(foldLine(line) + "\n")
	}
}

// writePhotoProperty writes photo property to the builder
func (v *VCard) writePhotoProperty(builder *strings.Builder) {
	if v.photo == "" {
		return
	}

	// Check if it's a URL or base64 data
	if strings.HasPrefix(v.photo, "http://") || strings.HasPrefix(v.photo, "https://") {
		// External URL
		line := fmt.Sprintf("PHOTO;VALUE=uri:%s", v.photo)
		builder.WriteString(foldLine(line) + "\n")
	} else if strings.HasPrefix(v.photo, "data:") {
		// Data URI (base64 encoded)
		line := fmt.Sprintf("PHOTO;ENCODING=b:%s", v.photo)
		builder.WriteString(foldLine(line) + "\n")
	} else {
		// Assume it's base64 data without data URI prefix
		line := fmt.Sprintf("PHOTO;ENCODING=b;TYPE=JPEG:%s", v.photo)
		builder.WriteString(foldLine(line) + "\n")
	}
}

// writeBirthdayProperty writes birthday property to the builder
func (v *VCard) writeBirthdayProperty(builder *strings.Builder) {
	if v.birthday == nil {
		return
	}

	// Format date according to vCard specification
	dateStr := v.birthday.Format("2006-01-02")
	line := fmt.Sprintf("BDAY:%s", dateStr)
	builder.WriteString(line + "\n")
}

// writeAnniversaryProperty writes anniversary property to the builder
func (v *VCard) writeAnniversaryProperty(builder *strings.Builder) {
	if v.anniversary == nil {
		return
	}

	// Anniversary is vCard 4.0 only
	if v.version == Version40 {
		dateStr := v.anniversary.Format("2006-01-02")
		line := fmt.Sprintf("ANNIVERSARY:%s", dateStr)
		builder.WriteString(line + "\n")
	}
}

// writeCustomProperties writes custom X- properties to the builder
func (v *VCard) writeCustomProperties(builder *strings.Builder) {
	for name, value := range v.customProps {
		if strings.HasPrefix(strings.ToUpper(name), "X-") && value != "" {
			line := fmt.Sprintf("%s:%s", strings.ToUpper(name), escapeValue(value))
			builder.WriteString(foldLine(line) + "\n")
		}
	}
}
